package folders

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FolderTransactionSuite struct {
	suite.Suite
	conn    *sql.DB
	mock    sqlmock.Sqlmock
	handler handler
	entity  *Folder
}

func (ts *FolderTransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.handler = handler{ts.conn}

	ts.entity = &Folder{
		Name:     "folder1",
	}
}

func (ts *FolderTransactionSuite) AfterTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(FolderTransactionSuite))
}
