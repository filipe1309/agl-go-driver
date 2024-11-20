package files

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestReadAllDB() {
	// Arrange
	setMockReadAllDB(ts.mock)

	// Act
	_, err := ReadAllDB(ts.conn, 0)

	// Assert
	assert.NoError(ts.T(), err)

	// if len(list) > 1 { // this doesn't work with sqlmock, because it doesn't filter the rows
	// 	t.Error("Expected 1 row, got", len(list))
	// }
}

func (ts *FileTransactionSuite) TestReadAllRootDB() {
	// Arrange
	setMockReadAllRootDB(ts.mock)

	// Act
	_, err := ReadAllRootDB(ts.conn)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadAllRootDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, nil, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)
}

func setMockReadAllDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(rows)
}
