package operator

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var fakeBucketName = "helix-saga"

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

func TestNewAliYunOss(t *testing.T) {
	type args struct {
		conf AliYunOssConfig
		ctx  context.Context
	}
	tests := []struct {
		name string
		args args
		want AliYunOss
	}{
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		fields  fields
		want    *oss.Client
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aliYunOss.Connector() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_aliYunOss_ListBuckets(t *testing.T) {
//	type fields struct {
//		mu   sync.RWMutex
//		conf AliYunOssConfig
//		ctx  context.Context
//	}
//	tests := []struct {
//		name      string
//		fields    fields
//		wantLsRes oss.ListBucketsResult
//		wantErr   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ays := &aliYunOss{
//				mu:   tt.fields.mu,
//				conf: tt.fields.conf,
//				ctx:  tt.fields.ctx,
//			}
//			gotLsRes, err := ays.ListBuckets()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("aliYunOss.ListBuckets() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotLsRes, tt.wantLsRes) {
//				t.Errorf("aliYunOss.ListBuckets() = %v, want %v", gotLsRes, tt.wantLsRes)
//			}
//		})
//	}
//}
//
//func Test_aliYunOss_CreateBucket(t *testing.T) {
//	type fields struct {
//		mu   sync.RWMutex
//		conf AliYunOssConfig
//		ctx  context.Context
//	}
//	type args struct {
//		bucketName string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ays := &aliYunOss{
//				mu:   tt.fields.mu,
//				conf: tt.fields.conf,
//				ctx:  tt.fields.ctx,
//			}
//			if err := ays.CreateBucket(tt.args.bucketName); (err != nil) != tt.wantErr {
//				t.Errorf("aliYunOss.CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func Test_aliYunOss_DeleteBucket(t *testing.T) {
//	type fields struct {
//		mu   sync.RWMutex
//		conf AliYunOssConfig
//		ctx  context.Context
//	}
//	type args struct {
//		bucketName string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ays := &aliYunOss{
//				mu:   tt.fields.mu,
//				conf: tt.fields.conf,
//				ctx:  tt.fields.ctx,
//			}
//			if err := ays.DeleteBucket(tt.args.bucketName); (err != nil) != tt.wantErr {
//				t.Errorf("aliYunOss.DeleteBucket() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func Test_aliYunOss_Bucket(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantB   *oss.Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			gotB, err := ays.Bucket(tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.Bucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("aliYunOss.Bucket() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func Test_aliYunOss_PutObject(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
		objectName string
		content    []byte
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
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			if err := ays.PutObject(tt.args.bucketName, tt.args.objectName, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.PutObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aliYunOss_GetObject(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
		objectName string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantContent []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			gotContent, err := ays.GetObject(tt.args.bucketName, tt.args.objectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("aliYunOss.GetObject() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func Test_aliYunOss_ListObjects(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantLsRes oss.ListObjectsResult
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			gotLsRes, err := ays.ListObjects(tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.ListObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLsRes, tt.wantLsRes) {
				t.Errorf("aliYunOss.ListObjects() = %v, want %v", gotLsRes, tt.wantLsRes)
			}
		})
	}
}

func Test_aliYunOss_DeleteObject(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
		objectName string
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
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			if err := ays.DeleteObject(tt.args.bucketName, tt.args.objectName); (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aliYunOss_GetEnvs(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			if gotRes := ays.GetEnvs(); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("aliYunOss.GetEnvs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_aliYunOss_CheckEnv(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		env string
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
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			if err := ays.CheckEnv(tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.CheckEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aliYunOss_GetContent(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
		env        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNc  NoticeContent
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			gotNc, err := ays.GetContent(tt.args.bucketName, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNc, tt.wantNc) {
				t.Errorf("aliYunOss.GetContent() = %v, want %v", gotNc, tt.wantNc)
			}
		})
	}
}

func Test_aliYunOss_UpdateContent(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		conf AliYunOssConfig
		ctx  context.Context
	}
	type args struct {
		bucketName string
		env        string
		nc         NoticeContent
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
			ays := &aliYunOss{
				mu:   tt.fields.mu,
				conf: tt.fields.conf,
				ctx:  tt.fields.ctx,
			}
			if err := ays.UpdateContent(tt.args.bucketName, tt.args.env, tt.args.nc); (err != nil) != tt.wantErr {
				t.Errorf("aliYunOss.UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
