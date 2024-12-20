package folders

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

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
		ID:   1,
		Name: "folder1",
	}
}

func (ts *FolderTransactionSuite) AfterTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func TestFolderSuite(t *testing.T) {
	suite.Run(t, new(FolderTransactionSuite))
}

func setMockFilesReadAllDB(mock sqlmock.Sqlmock, err bool) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test file 1", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test file 2", "testtype", "testpath", time.Now(), time.Now(), false)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(rows)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}

func setMockReadSubFolderDB(mock sqlmock.Sqlmock, err bool) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(2, 3, "Test sub name", time.Now(), time.Now(), false).
		AddRow(3, 3, "Test sub name 2", time.Now(), time.Now(), false)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`)).
		WithArgs(1).
		WillReturnRows(rows)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}

func setMockReadFolderDB(mock sqlmock.Sqlmock, err bool) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, 2, "Test folder 1", time.Now(), time.Now(), false)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}
