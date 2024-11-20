package folders

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestList() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	setMockReadRootSubFolderDB(ts.mock)
	setMockFilesReadAllRootDB(ts.mock)

	// Act
	ts.handler.List(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *FolderTransactionSuite) TestReadRootSubFolderDB() {
	// Arrange
	setMockReadRootSubFolderDB(ts.mock)

	// Act
	_, err := readRootSubFolderDB(ts.conn)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadRootSubFolderDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, "Test folder 1", time.Now(), time.Now(), false).
		AddRow(5, nil, "Test folder 1", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)
}

func setMockFilesReadAllRootDB(mock sqlmock.Sqlmock) {
	filesRow := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, nil, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(filesRow)
}
