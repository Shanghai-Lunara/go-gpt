package operator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"sync"
	"time"

	goftp "github.com/jlaffaye/ftp"
)

type FtpOperator interface {
	Conn() (c *goftp.ServerConn, err error)
	Quit(c *goftp.ServerConn)
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
	Name   string `json:"name"`
	Target string // target of symbolic link
	Type   int
	Size   uint64
	Time   time.Time `json:"time"`
}

func (f *ftp) Conn() (c *goftp.ServerConn, err error) {
	c, err = goftp.Dial(fmt.Sprintf("%s:%d", f.conf.Host, f.conf.Port), goftp.DialWithTimeout(time.Duration(f.conf.Timeout)*time.Second))
	if err != nil {
		return c, errors.New(fmt.Sprintf(errFtpConnDial, err))
	}
	defer f.Quit(c)
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

func (f *ftp) SetFileContent(res []byte) (err error) {

	return nil
}

func NewFtpOperator(c FtpConfig) FtpOperator {
	var f FtpOperator = &ftp{
		conf: c,
	}
	return f
}
