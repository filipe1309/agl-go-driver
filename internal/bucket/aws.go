package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSS3Config struct {
	Config         aws.Config
	BucketDownload string
	BucketUpload   string
}

func newAWSS3Session(cfg AWSS3Config) *awsSession {
	c, err := session.NewSession(&cfg.Config)
	if err != nil {
		panic(err)
	}
	return &awsSession{
		session:        c,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

type awsSession struct {
	session        *session.Session
	bucketDownload string
	bucketUpload   string
}

func (awsSession *awsSession) Download(src, dest string) (file *os.File, err error) {
	file, err = os.Create(dest)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(awsSession.session)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(awsSession.bucketDownload),
		Key:    aws.String(src),
	})
	return
}

func (awsSession *awsSession) Upload(file io.Reader, key string) error {
	return nil
}

func (awsSession *awsSession) Delete(key string) error {
	return nil
}
