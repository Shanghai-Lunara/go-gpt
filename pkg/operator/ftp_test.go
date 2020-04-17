package operator

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	goftp "github.com/jlaffaye/ftp"
)

var fakeFtpConfig = FtpConfig{
	Username: "anonymous",
	Password: "anonymous",
	WorkDir:  "/",
	Host:     "127.0.0.1",
	Port:     21,
}

func Test_ftp_Conn(t *testing.T) {
	mock, c := openConn(t, fmt.Sprintf("%s:%d", fakeFtpConfig.Host, fakeFtpConfig.Port), goftp.DialWithTimeout(5*time.Second))
	_ = mock
	_ = c
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantC   *goftp.ServerConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ftp{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
			}
			gotC, err := f.Conn()
			if (err != nil) != tt.wantErr {
				t.Errorf("ftp.Conn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("ftp.Conn() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func Test_ftp_Quit(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	type args struct {
		c *goftp.ServerConn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ftp{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
			}
			f.Quit(tt.args.c)
		})
	}
}

