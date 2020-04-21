package operator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"k8s.io/klog"
)

type AliYunOss interface {
	Connector() (*oss.Client, error)
	ListBuckets() (lsRes oss.ListBucketsResult, err error)
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	Bucket(bucketName string) (b *oss.Bucket, err error)
	PutObject(bucketName, objectName string, content []byte) error
	GetObject(bucketName, objectName string) (content []byte, err error)
	ListObjects(bucketName string) (lsRes oss.ListObjectsResult, err error)
	DeleteObject(bucketName, objectName string) error
	GetEnvs() (res map[string]string)
	GetContent(bucketName, env string) (nc NoticeContent, err error)
	UpdateContent(bucketName, env string, nc NoticeContent) error
}

const (
	errEnvWasNotExisted = "the env name:%s was not existed"
)

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
	err = client.CreateBucket(bucketName, oss.ACL(oss.ACLPublicReadWrite))
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

func (ays *aliYunOss) PutObject(bucketName, objectName string, content []byte) error {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = b.PutObject(objectName, strings.NewReader(string(content)), oss.ACL(oss.ACLPublicReadWrite))
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (ays *aliYunOss) GetObject(bucketName, objectName string) (content []byte, err error) {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return content, err
	}
	body, err := b.GetObject(objectName, oss.ACL(oss.ACLPublicReadWrite))
	content, err = ioutil.ReadAll(body)
	if err = body.Close(); err != nil {
		klog.V(2).Info(err)
	}
	return content, err
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

func (ays *aliYunOss) DeleteObject(bucketName, objectName string) error {
	b, err := ays.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = b.DeleteObject(objectName, oss.ACL(oss.ACLPublicReadWrite))
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (ays *aliYunOss) GetEnvs() (res map[string]string) {
	res = make(map[string]string, 0)
	for _, v := range ays.conf.Envs {
		res[v.Name] = fmt.Sprintf("%s/dev/%s.html", ays.conf.ProxyUrl, v.Value)
	}
	return res
}

func (ays *aliYunOss) CheckEnv(env string) error {
	match := false
	for _, v := range ays.conf.Envs {
		if v.Name == env {
			match = true
		}
	}
	if match == true {
		return nil
	}
	err := errors.New(fmt.Sprintf(errEnvWasNotExisted, env))
	klog.V(2).Info(err)
	return err
}

func (ays *aliYunOss) GetContent(bucketName, env string) (nc NoticeContent, err error) {
	if err = ays.CheckEnv(env); err != nil {
		return nc, err
	}
	if bucketName == "" {
		bucketName = ays.conf.BucketName
	}
	tmp, err := ays.GetObject(bucketName, fmt.Sprintf("%s.json", env))
	if err != nil {
		return nc, nil
	}
	err = json.Unmarshal(tmp, &nc)
	if err != nil {
		klog.V(2).Info(err)
	}
	return nc, nil
}

func (ays *aliYunOss) UpdateContent(bucketName, env string, nc NoticeContent) error {
	if err := ays.CheckEnv(env); err != nil {
		return err
	}
	if bucketName == "" {
		bucketName = ays.conf.BucketName
	}
	// update {env}.json
	data, err := json.Marshal(nc)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	err = ays.PutObject(bucketName, fmt.Sprintf("%s.json", env), []byte(data))
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	// update deb/{env}.html
	tmp, err := ays.GetObject(bucketName, "index.html")
	if err != nil {
		return err
	}
	str := strings.Replace(string(tmp), "{{title}}", nc.Title, -1)
	str = strings.Replace(str, "{{time}}", nc.Time, -1)
	str = strings.Replace(str, "{{content}}", nc.Content, -1)
	err = ays.PutObject(bucketName, fmt.Sprintf("dev/%s.html", env), []byte(str))
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}
