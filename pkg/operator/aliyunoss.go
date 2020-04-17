package operator

import (
	"context"
	"k8s.io/klog"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYunOss interface {
	Connector() (*oss.Client, error)
	ListBuckets() (lsRes oss.ListBucketsResult, err error)
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	Bucket(bucketName string) (b *oss.Bucket, err error)
	PutObject(bucketName, objectName, fileDir string) error
	ListObjects(bucketName string) (lsRes oss.ListObjectsResult, err error)
	DeleteObjects(bucketName, objectName string) error
}

func NewAliYunOss(conf AliYunOssConfig, ctx context.Context) AliYunOss {
	var ays AliYunOss = &aliYunOss{
		conf: conf,
		ctx:  ctx,
	}
	return ays
}

type aliYunOss struct {
	mu   sync.RWMutex
	conf AliYunOssConfig
	ctx  context.Context
}

func (ays *aliYunOss) Connector() (*oss.Client, error) {
	client, err := oss.New(ays.conf.EndPoint, ays.conf.AccessKeyID, ays.conf.AccessKeySecret)
	if err != nil {
		klog.V(2).Info(err)
	}
	return client, err
}

func (ays *aliYunOss) ListBuckets() (lsRes oss.ListBucketsResult, err error) {
	client, err := ays.Connector()
	if err != nil {
		return lsRes, err
	}
	if lsRes, err = client.ListBuckets(); err != nil {
		klog.V(2).Info(err)
	}
	return lsRes, err
}

func (ays *aliYunOss) CreateBucket(bucketName string) error {
	client, err := ays.Connector()
	if err != nil {
		return err
	}
	err = client.CreateBucket(bucketName)
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (ays *aliYunOss) DeleteBucket(bucketName string) error {
	client, err := ays.Connector()
	if err != nil {
		return err
	}
	err = client.DeleteBucket(bucketName)
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (ays *aliYunOss) Bucket(bucketName string) (b *oss.Bucket, err error) {
	client, err := ays.Connector()
	if err != nil {
		return b, err
	}
	b, err = client.Bucket(bucketName)
	if err != nil {
		klog.Info(err)
	}
	return b, err
}

func (ays *aliYunOss) PutObject(bucketName, objectName, fileDir string) error {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = b.PutObjectFromFile(objectName, fileDir)
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (ays *aliYunOss) ListObjects(bucketName string) (lsRes oss.ListObjectsResult, err error) {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return lsRes, err
	}
	lsRes, err = b.ListObjects()
	if err != nil {
		klog.V(2).Info(err)
	}
	return lsRes, err
}

func (ays *aliYunOss) DeleteObjects(bucketName, objectName string) error {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = b.DeleteObject(objectName)
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}
