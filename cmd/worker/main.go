package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/filipe1309/agl-go-driver/internal/bucket"
	"github.com/filipe1309/agl-go-driver/internal/queue"
)

func main() {
	// Queue config
	queueConfig := queue.RabbitMQConfig{
		URL:       "amqp://" + os.Getenv("QUEUE_RABBITMQ_URL"),
		TopicName: os.Getenv("QUEUE_RABBITMQ_TOPIC"),
		Timeout:   time.Second * 30,
	}

	queueConn, err := queue.New(queue.RabbitMQ, queueConfig)
	if err != nil {
		panic(err)
	}

	// create a channel to receive messages
	queueChannel := make(chan queue.QueueDTO)
	go queueConn.Consume(queueChannel)

	// Bucket config
	bucketConfig := bucket.AWSS3Config{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), os.Getenv("AWS_SESSION_TOKEN")),
		},
		BucketDownload: os.Getenv("WORKER_BUCKET_AWS_S3_DOWNLOAD"),
		BucketUpload:   os.Getenv("WORKER_BUCKET_AWS_S3_UPLOAD"),
	}

	bucket, err := bucket.New(bucket.AWSS3Provider, bucketConfig)
	if err != nil {
		panic(err)
	}

	log.Println("waiting for messages...")
	for msg := range queueChannel {
		src := msg.Path
		dest := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)

		log.Printf("downloading file %s to %s", src, dest)

		err := bucket.Download(src, dest)
		if err != nil {
			log.Printf("error downloading file: %s", err)
			continue
		}

		file, err := os.Open(dest)
		if err != nil {
			log.Printf("error opening file: %s", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("error reading file: %s", err)
			continue
		}

		var buffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&buffer)
		_, err = gzipWriter.Write(body)
		if err != nil {
			log.Printf("error compressing file: %s", err)
			continue
		}

		if err := gzipWriter.Close(); err != nil {
			log.Printf("error closing gzip writer: %s", err)
			continue
		}

		gzipReader, err := gzip.NewReader(&buffer)
		if err != nil {
			log.Printf("error creating gzip reader: %s", err)
			continue
		}

		if err := bucket.Upload(gzipReader, src); err != nil {
			log.Printf("error uploading file: %s", err)
			continue
		}

		if err := gzipReader.Close(); err != nil {
			log.Printf("error closing gzip reader: %s", err)
			continue
		}

		if err := os.Remove(dest); err != nil {
			log.Printf("error removing file: %s", err)
			continue
		}

		log.Printf("file %s processed with success", src)
	}
}
