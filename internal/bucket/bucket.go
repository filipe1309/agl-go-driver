package bucket

import (
	"fmt"
	"io"
	"reflect"
)

const (
	AWSS3Provider BucketType = iota
	MockBucketProvider
)

type BucketType int

func New(bt BucketType, cfg any) (b *Bucket, err error) {
	b = new(Bucket)

	rt := reflect.TypeOf(cfg)

	switch bt {
	case AWSS3Provider:
		if rt.Name() != "AWSS3Config" {
			return nil, fmt.Errorf("invalid aws s3 config type: %s", rt.Name())
		}

		b.provider = newAWSS3Session(cfg.(AWSS3Config))
	case MockBucketProvider:
		b.provider = &MockBucket{
			content: make(map[string][]byte),
		}
	default:
		return nil, fmt.Errorf("type not supported")
	}
	return
}

type BucketInterface interface {
	Upload(file io.Reader, key string) error
	Download(src string, dest string) error
	Delete(src string) error
}

type Bucket struct {
	provider BucketInterface
}

func (b *Bucket) Upload(file io.Reader, key string) error {
	return b.provider.Upload(file, key)
}

func (b *Bucket) Download(src, dest string) error {
	return b.provider.Download(src, dest)
}

func (b *Bucket) Delete(key string) error {
	return b.provider.Delete(key)
}
