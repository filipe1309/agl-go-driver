package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/filipe1309/agl-go-driver/internal/bucket"
	"github.com/filipe1309/agl-go-driver/internal/files"
	"github.com/filipe1309/agl-go-driver/internal/folders"
	"github.com/filipe1309/agl-go-driver/internal/queue"
	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/pkg/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, bucket, queueConn := getSeessions()

	// Define endpoints
	r := chi.NewRouter()
	r.Post("/auth", auth.NewHandlerAuth(func(login, password string) (auth.Authenticated, error) {
		return users.Authenticate(login, password)
	}))

	files.SetRoutes(r, db, bucket, queueConn)
	folders.SetRoutes(r, db)
	users.SetRoutes(r, db)

	// Start server
	log.Println("Server running on port 8080")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), r)
}

func getSeessions() (*sql.DB, *bucket.Bucket, *queue.Queue) {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Queue config
	queueConfig := queue.RabbitMQConfig{
		URL:       os.Getenv("QUEUE_RABBITMQ_URL"),
		TopicName: os.Getenv("QUEUE_RABBITMQ_TOPIC"),
		Timeout:   time.Second * 30,
	}

	queueConn, err := queue.New(queue.RabbitMQ, queueConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Bucket config
	bucketConfig := bucket.AWSS3Config{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: os.Getenv("BUCKET_AWS_S3_DOWNLOAD"),
		BucketUpload:   os.Getenv("BUCKET_AWS_S3_UPLOAD"),
	}

	bucket, err := bucket.New(bucket.AWSS3Provider, bucketConfig)
	if err != nil {
		log.Fatal(err)
	}

	return db, bucket, queueConn
}
