package folders

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestList() {
	tcs := []struct {
		Name                          string
		MockReadAllDBWithErr          bool
		MockFilesReadAllRootDBWithErr bool
		ExpectedStatusCode            int
	}{
		{
			Name:                          "Success",
			MockReadAllDBWithErr:          false,
			MockFilesReadAllRootDBWithErr: false,
			ExpectedStatusCode:            http.StatusOK,
		},
		{
			Name:                          "DB error - read root sub folder",
			MockReadAllDBWithErr:          true,
			MockFilesReadAllRootDBWithErr: false,
			ExpectedStatusCode:            http.StatusInternalServerError,
		},
		{
			Name:                          "DB error - read files",
			MockReadAllDBWithErr:          false,
			MockFilesReadAllRootDBWithErr: true,
			ExpectedStatusCode:            http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)

		setMockReadRootSubFolderDB(ts.mock, tc.MockReadAllDBWithErr)
		if !tc.MockReadAllDBWithErr {
			setMockFilesReadAllRootDB(ts.mock, tc.MockFilesReadAllRootDBWithErr)
		}

		// Act
		ts.handler.List(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FolderTransactionSuite) TestReadRootSubFolderDB() {
	// Arrange
	setMockReadRootSubFolderDB(ts.mock, false)

	// Act
	_, err := readRootSubFolderDB(ts.conn)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadRootSubFolderDB(mock sqlmock.Sqlmock, err bool) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, "Test folder 1", time.Now(), time.Now(), false).
		AddRow(5, nil, "Test folder 1", time.Now(), time.Now(), false)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}

func setMockFilesReadAllRootDB(mock sqlmock.Sqlmock, err bool) {
	filesRow := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, nil, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(filesRow)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}
