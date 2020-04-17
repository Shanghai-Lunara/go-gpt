package operator

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func GetConf() (ayoc AliYunOssConfig, err error) {
	ayoc = AliYunOssConfig{
		EndPoint:        "121",
		Bucket:          "2121",
		AccessKeyID:     "2121",
		AccessKeySecret: "212121",
		FileDirectory:   "/tmp/",
		Envs: []AliYunOssEnv{
			{
				Name:  "dev",
				Value: "dev",
			},
			{
				Name:  "test",
				Value: "test",
			},
			{
				Name:  "product",
				Value: "product",
			},
		},
	}
	return ayoc, nil
}

func GetFakeErrConf() AliYunOssConfig {
	return AliYunOssConfig{
		EndPoint:        "121",
		Bucket:          "2121",
		AccessKeyID:     "2121",
		AccessKeySecret: "212121",
		FileDirectory:   "/tmp/",
		Envs: []AliYunOssEnv{
			{
				Name:  "dev",
				Value: "dev",
			},
			{
				Name:  "test",
				Value: "test",
			},
			{
				Name:  "product",
				Value: "product",
			},
		},
	}
}

func GetArgs() (conf AliYunOssConfig, ctx context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf, err := GetConf()
	if err != nil {
		panic(err)
	}
	return conf, ctx
}

func TestNewAliYunOss(t *testing.T) {
	type args struct {
		conf AliYunOssConfig
		ctx  context.Context
	}
	conf, ctx := GetArgs()
	tmp := args{
		conf: conf,
		ctx:  ctx,
	}
	var ays AliYunOss = &aliYunOss{
		conf: conf,
		ctx:  ctx,
	}
	tests := []struct {
		name string
		args args
		want AliYunOss
	}{
		{
			name: "case_1",
			args: tmp,
			want: ays,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAliYunOss(tt.args.conf, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAliYunOss() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aliYunOss_Connector(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	conf, ctx := GetArgs()
	tests := []struct {
		name    string
		fields  fields
		want    *oss.Client
		wantErr bool
	}{
		{
			name: "case_1",
			fields: fields{
				conf: conf,
				ctx:  ctx,
			},
			wantErr: false,
		},
		{
			name: "case_2",
			fields: fields{
				conf: GetFakeErrConf(),
				ctx:  ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			got, err := ays.Connector()
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.Connector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("aliYunOss.Connector() = %v, want %v", got, tt.want)
			//}
			_ = got
		})
	}
}

func Test_aliYunOss_ListBuckets(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	conf, ctx := GetArgs()
	_ = conf
	tests := []struct {
		name      string
		fields    fields
		wantLsRes oss.ListBucketsResult
		wantErr   bool
	}{
		{
			name: "case_1",
			fields: fields{
				conf: conf,
				ctx:  ctx,
			},
			wantErr: false,
		},
		{
			name: "case_2",
			fields: fields{
				conf: GetFakeErrConf(),
				ctx:  ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			gotLsRes, err := ays.ListBuckets()
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = gotLsRes
			//fmt.Print(gotLsRes)
			//if !reflect.DeepEqual(gotLsRes, tt.wantLsRes) {
			//	t.Errorf("aliYunOss.List() = %v, want %v", gotLsRes, tt.wantLsRes)
			//}
		})
	}
}
