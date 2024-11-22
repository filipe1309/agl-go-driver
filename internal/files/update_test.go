package files

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestUpdate() {
	tcs := []struct {
		Name                string
		IDStr               string
		WithMock            bool
		MockFile            *File
		MockReadDBWithErr   bool
		MockUpdateDBWithErr bool
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success",
			IDStr:               "1",
			WithMock:            true,
			MockFile:            &File{ID: 1, Name: "Test file 2"},
			MockReadDBWithErr:   false,
			MockUpdateDBWithErr: false,
			ExpectedStatusCode:  http.StatusOK,
		},
		{
			Name:                "Invalid file - no name",
			IDStr:               "1",
			WithMock:            true,
			MockFile:            &File{ID: 1},
			MockReadDBWithErr:   false,
			MockUpdateDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
		{
			Name:                "Invalid json - empty body",
			IDStr:               "1",
			WithMock:            true,
			MockFile:            nil,
			MockReadDBWithErr:   false,
			MockUpdateDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
		{
			Name:                "Invalid url param - id",
			IDStr:               "A",
			WithMock:            false,
			MockFile:            &File{ID: -1, Name: "Test file 1"},
			MockReadDBWithErr:   false,
			MockUpdateDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
		{
			Name:                "DB error - read",
			IDStr:               "1",
			WithMock:            true,
			MockFile:            &File{ID: 1, Name: "Test file 1"},
			MockReadDBWithErr:   true,
			MockUpdateDBWithErr: false,
			ExpectedStatusCode:  http.StatusInternalServerError,
		},

		{
			Name:                "DB error - update",
			IDStr:               "1",
			WithMock:            true,
			MockFile:            &File{ID: 1, Name: "Test file 2"},
			MockReadDBWithErr:   false,
			MockUpdateDBWithErr: true,
			ExpectedStatusCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		var b bytes.Buffer
		if tc.MockFile != nil {
			err := json.NewEncoder(&b).Encode(tc.MockFile)
			assert.NoError(ts.T(), err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/files/{id}", &b)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockReadDB(ts.mock, tc.MockReadDBWithErr)
			if !tc.MockReadDBWithErr && tc.MockFile != nil && tc.MockFile.Name != "" {
				setMockUpdateDB(ts.mock, 1, tc.MockUpdateDBWithErr)
			}
		}

		// Act
		ts.handler.Update(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FileTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1, false)

	// Act
	_, err := UpdateDB(ts.conn, 1, &File{
		Name: "Test file 2",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs("Test file 2", sqlmock.AnyArg(), false, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err {
		exp.WillReturnError(sqlmock.ErrCancelled)
	}
}
