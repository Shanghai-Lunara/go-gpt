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
	jsoniter "github.com/json-iterator/go"
	"k8s.io/klog"
)

type FtpOperator interface {
	Conn() (c *goftp.ServerConn, err error)
	Quit(c *goftp.ServerConn)
	List() (res []Entry, err error)
	ReadFileContent(fileName string) (res []byte, err error)
	SetFileContent(fileName string, content []byte) (err error)
	UploadFile(sourcePath, fileName string) (err error)
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

func (f *ftp) List() (res []Entry, err error) {
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

func (f *ftp) SetFileContent(fileName string, content []byte) (err error) {
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
	file, err := os.Open(fmt.Sprintf("%s/%s", sourcePath, fileName))
	if err != nil {
		klog.V(2).Info(err)
	}
	if err = c.Stor(fmt.Sprintf("%s/%s", f.conf.WorkDir, fileName), file); err != nil {
		return err
	}
	return nil
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewFtpOperator(c FtpConfig) FtpOperator {
	var f FtpOperator = &ftp{
		conf: c,
	}
	if res, err := f.ReadFileContent("introduce_2020022602.txt"); err != nil {
		klog.Info(err)
	} else {
		klog.Info(string(res))
	}
	if res, err := f.List(); err != nil {
		klog.Infof("List err:%v", err)
	} else {
		t, err := json.Marshal(res)
		if err != nil {
			klog.Info(err)
		}
		klog.Info("t:", string(t))
		GetVersionByFilter(res, "introduce")
	}
	if err := f.SetFileContent("abc.txt", []byte("test\nabc1213")); err != nil {
		klog.Info(err)
	}
	return f
}

func GetVersionByFilter(source []Entry, specNamePrefix string) (version int) {
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
		t := `introduce_([\d]{4})([\d]{4})([\d]{2}).txt`
		re := regexp.MustCompile(t)
		res := re.FindStringSubmatch(v.Name)
		_, err = strconv.Atoi(res[3])
		if err != nil {
			klog.V(2).Info(err)
			version = 0
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

func GetTodayTimePrefix() string {
	//t := time.Now().Format("20060102")
	return time.Now().Format("20060102")
}