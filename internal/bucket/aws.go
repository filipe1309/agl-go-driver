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
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

func newAWSS3Session(cfg AWSS3Config) *awsSession {
	c, err := session.NewSession(cfg.Config)
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

func (awsSession *awsSession) Download(src, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(awsSession.session)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(awsSession.bucketDownload),
		Key:    aws.String(src),
	})
	return err
}

func (awsSession *awsSession) Upload(file io.Reader, key string) error {
	uploader := s3manager.NewUploader(awsSession.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsSession.bucketUpload),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (awsSession *awsSession) Delete(key string) error {
	svc := s3.New(awsSession.session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsSession.bucketUpload),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsSession.bucketUpload),
		Key:    aws.String(key),
	})
}
