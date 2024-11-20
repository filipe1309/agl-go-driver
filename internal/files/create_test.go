package files

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/common"
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestCreate() {
	// Arrange
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	file, err := os.Open("./testdata/test-image-1.jpg")
	assert.NoError(ts.T(), err)

	fw, err := mw.CreateFormFile("file", "test-image-1.jpg")
	assert.NoError(ts.T(), err)

	_, err = io.Copy(fw, file)
	assert.NoError(ts.T(), err)

	mw.Close()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/files", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	setMockInsertDB(ts.mock, ts.entity, nil)

	// Act
	ts.handler.Create(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusCreated, rr.Code)
}

func (ts *FileTransactionSuite) TestInsertRootDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity, nil)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FileTransactionSuite) TestInsertWithFolderDB() {
	// Arrange
	ts.entity.FolderID = common.ValidNullInt64(2)
	setMockInsertDB(ts.mock, ts.entity, 2)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *File, folderID any) {
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(folderID, 1, entity.Name, entity.Type, entity.Path, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
