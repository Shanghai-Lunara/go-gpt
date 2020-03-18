package operator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	goftp "github.com/jlaffaye/ftp"
	"k8s.io/klog"
)

type FtpOperator interface {
	Conn() (c *goftp.ServerConn, err error)
	Quit(c *goftp.ServerConn)
	List(filter string) (res []Entry, err error)
	ReadFileContent(fileName string) (res []byte, err error)
	WriteFileContent(fileName string, content []byte) (err error)
	UploadFile(sourcePath, fileName string) (err error)
	GetNextVersion() (version string, err error)
}

type ftp struct {
	mu sync.RWMutex

	conf FtpConfig
}

const (
	errFtpConnDial = "ftp dial err:%v"
	errQuitErr     = "ftp quit err:%v"
	errFtpLogin    = "ftp login err:%v"
)

// Entry describes a file and is returned by List().
type Entry struct {
	Name   string    `json:"name"`
	Target string    `json:"target"` // target of symbolic link
	Type   int       `json:"type"`
	Size   uint64    `json:"size"`
	Time   time.Time `json:"time"`
}

func (f *ftp) Conn() (c *goftp.ServerConn, err error) {
	c, err = goftp.Dial(fmt.Sprintf("%s:%d", f.conf.Host, f.conf.Port), goftp.DialWithTimeout(time.Duration(f.conf.Timeout)*time.Second))
	if err != nil {
		return c, errors.New(fmt.Sprintf(errFtpConnDial, err))
	}
	err = c.Login(f.conf.Username, f.conf.Password)
	if err != nil {
		return c, errors.New(fmt.Sprintf(errFtpLogin, err))
	}
	return c, nil
}

func (f *ftp) Quit(c *goftp.ServerConn) {
	if err := c.Quit(); err != nil {
		klog.V(2).Info(errQuitErr, err)
	}
}

func (f *ftp) List(filter string) (res []Entry, err error) {
	c, err := f.Conn()
	if err != nil {
		return res, err
	}
	defer f.Quit(c)
	ret, err := c.List(f.conf.WorkDir)
	if err != nil {
		return res, err
	}
	res = make([]Entry, 0)
	for _, v := range ret {
		if filter != "" {
			matched, err := regexp.Match(filter, []byte(v.Name))
			if err != nil {
				klog.V(2).Info(err)
			}
			if !matched {
				continue
			}
		}
		t := Entry{
			Name:   v.Name,
			Target: v.Target,
			Type:   int(v.Type),
			Size:   v.Size,
			Time:   v.Time,
		}
		res = append(res, t)
	}
	return res, nil
}

func (f *ftp) ReadFileContent(fileName string) (res []byte, err error) {
	c, err := f.Conn()
	if err != nil {
		return res, err
	}
	defer f.Quit(c)
	cont, err := c.Retr(fmt.Sprintf("%s/%s", f.conf.WorkDir, fileName))
	if err != nil {
		return res, err
	}
	res, err = ioutil.ReadAll(cont)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (f *ftp) WriteFileContent(fileName string, content []byte) (err error) {
	c, err := f.Conn()
	if err != nil {
		return err
	}
	defer f.Quit(c)
	err = c.Stor(fmt.Sprintf("%s/%s", f.conf.WorkDir, fileName), bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	return nil
}

func (f *ftp) UploadFile(sourcePath, fileName string) (err error) {
	c, err := f.Conn()
	if err != nil {
		return err
	}
	defer f.Quit(c)
	file, err := os.Open(sourcePath)
	if err != nil {
		klog.V(2).Info(err)
	}
	if err = c.Stor(fmt.Sprintf("%s/%s", f.conf.WorkDir, fileName), file); err != nil {
		return err
	}
	return nil
}

func (f *ftp) GetNextVersion() (version string, err error) {
	res, err := f.List("")
	if err != nil {
		return GetNextVersionString(0), errors.New(fmt.Sprintf("List err:%v", err))
	}
	v := GetTodayVersionByFilter(res, "introduce", "txt")
	return GetNextVersionString(v), nil
}

func NewFtpOperator(c FtpConfig) FtpOperator {
	var f FtpOperator = &ftp{
		conf: c,
	}
	return f
}

func GetTodayVersionByFilter(source []Entry, specNamePrefix, specNameSuffix string) (version int) {
	version = 0
	for _, v := range source {
		matched, err := regexp.Match(specNamePrefix, []byte(v.Name))
		if err != nil {
			klog.V(2).Info(err)
			continue
		}
		if matched == false {
			continue
		}
		re := regexp.MustCompile(fmt.Sprintf(`%s_%s([\d]{2}).%s`, specNamePrefix, time.Now().Format("20060102"), specNameSuffix))
		res := re.FindStringSubmatch(v.Name)
		if len(res) < 2 {
			continue
		}
		tmp, err := strconv.Atoi(res[1])
		if err != nil {
			klog.V(2).Info(err)
		}
		if tmp > version {
			version = tmp
		}
		klog.Infof("Filter version:%v", fmt.Sprintf("%d", version))
	}
	return version
}

func GetNextVersionString(version int) (str string) {
	version++
	if version >= 10 {
		return strconv.Itoa(version)
	}
	return fmt.Sprintf("0%d", version)
}
