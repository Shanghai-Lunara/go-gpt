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
	mock, c := openConn(t, fmt.Sprintf("%s%d", fakeFtpConfig.Host, fakeFtpConfig.Port), goftp.DialWithTimeout(5*time.Second))
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

func Test_ftp_List(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []Entry
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
			gotRes, err := f.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("ftp.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ftp.List() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_ftp_ReadFileContent(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []byte
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
			gotRes, err := f.ReadFileContent(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ftp.ReadFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ftp.ReadFileContent() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_ftp_SetFileContent(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	type args struct {
		fileName string
		content  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := f.SetFileContent(tt.args.fileName, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("ftp.SetFileContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewFtpOperator(t *testing.T) {
	type args struct {
		c FtpConfig
	}
	tests := []struct {
		name string
		args args
		want FtpOperator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFtpOperator(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFtpOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ftp_UploadFile(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf FtpConfig
	}
	type args struct {
		sourcePath string
		fileName   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := f.UploadFile(tt.args.sourcePath, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("ftp.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetVersionByFilter(t *testing.T) {
	type args struct {
		source         []Entry
		specNamePrefix string
	}
	tests := []struct {
		name        string
		args        args
		wantVersion int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVersion := GetVersionByFilter(tt.args.source, tt.args.specNamePrefix); gotVersion != tt.wantVersion {
				t.Errorf("GetVersionByFilter() = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestGetNextVersionString(t *testing.T) {
	type args struct {
		version int
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := GetNextVersionString(tt.args.version); gotStr != tt.wantStr {
				t.Errorf("GetNextVersionString() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestGetTodayTimePrefix(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "case1",
			want: "20200313",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTodayTimePrefix(); got != tt.want {
				t.Errorf("GetTodayTimePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
