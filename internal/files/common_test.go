package files

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/bucket"
	"github.com/filipe1309/agl-go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileTransactionSuite struct {
	suite.Suite
	conn    *sql.DB
	mock    sqlmock.Sqlmock
	handler handler
	entity  *File
}

func (ts *FileTransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	mockBucket, err := bucket.New(bucket.MockBucketProvider, nil)
	assert.NoError(ts.T(), err)

	mockQueue, err := queue.New(queue.MockQueueProvider, nil)
	assert.NoError(ts.T(), err)

	ts.handler = handler{ts.conn, mockBucket, mockQueue}

	ts.entity, err = New(1, "test-image-1.jpg", "application/octet-stream", "/test-image-1.jpg")
	assert.NoError(ts.T(), err)
}

func (ts *FileTransactionSuite) AfterTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileTransactionSuite))
}
