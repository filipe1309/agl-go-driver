package files

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestReadDB() {
	// Arrange
	setMockReadDB(ts.mock)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)
}
