package files

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func (ts *FileTransactionSuite) TestCreate() {
	tcs := []struct {
		Name                string
		FolderID            any
		WithUpload          bool
		MockFile            *File
		WithMockDB          bool
		MockInsertDBWithErr bool
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success - insert on root",
			FolderID:            nil,
			WithUpload:          true,
			MockFile:            &File{ID: 1, Name: "test-image-1.jpg"},
			WithMockDB:          true,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusCreated,
		},
		// {
		// 	Name:                "Success - insert on existent folder",
		// 	FolderID:            int64(1),
		// 	WithUpload:          true,
		// 	MockFile:            &File{ID: 1, Name: "test-image-1.jpg"},
		// 	WithMockDB:          true,
		// 	MockInsertDBWithErr: false,
		// 	ExpectedStatusCode:  http.StatusCreated,
		// },
		// {
		// 	Name:                "DB error - insert",
		// 	FolderID:            int64(1),
		// 	WithUpload:          true,
		// 	MockFile:            &File{ID: 1, Name: "test-image-1.jpg"},
		// 	WithMockDB:          true,
		// 	MockInsertDBWithErr: true,
		// 	ExpectedStatusCode:  http.StatusInternalServerError,
		// },
		// {
		// 	Name:                "Invalid param - folder id",
		// 	FolderID:            "A",
		// 	WithUpload:          true,
		// 	MockFile:            &File{ID: 1, Name: "test-image-1.jpg"},
		// 	WithMockDB:          false,
		// 	MockInsertDBWithErr: false,
		// 	ExpectedStatusCode:  http.StatusBadRequest,
		// },
		// {
		// 	Name:                "No file",
		// 	FolderID:            int64(1),
		// 	WithUpload:          false,
		// 	MockFile:            nil,
		// 	WithMockDB:          false,
		// 	MockInsertDBWithErr: false,
		// 	ExpectedStatusCode:  http.StatusBadRequest,
		// },
	}

	for _, tc := range tcs {
		// Arrange
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		if tc.WithUpload {
			file, err := os.Open(fmt.Sprintf("./testdata/%s", tc.MockFile.Name))
			assert.NoError(ts.T(), err)

			fw, err := mw.CreateFormFile("file", tc.MockFile.Name)
			assert.NoError(ts.T(), err)

			_, err = io.Copy(fw, file)
			assert.NoError(ts.T(), err)

			mw.Close()
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/files", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		req = req.WithContext(common.SetUserIDInContext(req.Context(), 1))

		if tc.FolderID != nil {
			req.Form = make(map[string][]string)
			req.Form.Add("folder_id", fmt.Sprintf("%v", tc.FolderID))
			// w, err := mw.CreateFormField("folder_id")
			// assert.NoError(ts.T(), err)
			// w.Write([]byte(fmt.Sprintf("%d", tc.FolderID)))
		}

		if tc.WithMockDB {
			setMockInsertDB(ts.mock, ts.entity, tc.FolderID, tc.MockInsertDBWithErr)
		}

		// Act
		ts.handler.Create(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FileTransactionSuite) TestInsertRootDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity, nil, false)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FileTransactionSuite) TestInsertWithFolderDB() {
	// Arrange
	ts.entity.FolderID = null.IntFrom(2)
	setMockInsertDB(ts.mock, ts.entity, 2, false)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *File, folderID any, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)).
		WithArgs(folderID, 1, entity.Name, entity.Type, entity.Path, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	if err {
		exp.WillReturnError(sqlmock.ErrCancelled)
	}
}
