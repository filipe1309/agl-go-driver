package folders

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestGetByID() {
	tcs := []struct {
		Name                       string
		ID                         int
		IDStr                      string
		WithMockDB                 bool
		MockReadFolderDBWithErr    bool
		MockReadSubFolderDBWithErr bool
		MockFilesReadAllDBWithErr  bool
		ExpectedStatusCode         int
	}{
		{
			Name:                       "Success",
			ID:                         1,
			IDStr:                      "1",
			WithMockDB:                 true,
			MockReadFolderDBWithErr:    false,
			MockReadSubFolderDBWithErr: false,
			MockFilesReadAllDBWithErr:  false,
			ExpectedStatusCode:         http.StatusOK,
		},
		{
			Name:                      "Invalid url param - id",
			ID:                        0,
			IDStr:                     "A",
			WithMockDB:                false,
			MockReadFolderDBWithErr:   false,
			MockFilesReadAllDBWithErr: false,
			ExpectedStatusCode:        http.StatusBadRequest,
		},
		{
			Name:                       "DB error - read folder",
			ID:                         1,
			IDStr:                      "1",
			WithMockDB:                 true,
			MockReadFolderDBWithErr:    true,
			MockReadSubFolderDBWithErr: false,
			MockFilesReadAllDBWithErr:  false,
			ExpectedStatusCode:         http.StatusInternalServerError,
		},
		{
			Name:                       "DB error - read sub folder",
			ID:                         1,
			IDStr:                      "1",
			WithMockDB:                 true,
			MockReadFolderDBWithErr:    false,
			MockReadSubFolderDBWithErr: true,
			MockFilesReadAllDBWithErr:  false,
			ExpectedStatusCode:         http.StatusInternalServerError,
		},
		{
			Name:                       "DB error - read all files",
			ID:                         1,
			IDStr:                      "1",
			WithMockDB:                 true,
			MockReadFolderDBWithErr:    false,
			MockReadSubFolderDBWithErr: false,
			MockFilesReadAllDBWithErr:  true,
			ExpectedStatusCode:         http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/folders/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMockDB {
			setMockReadFolderDB(ts.mock, tc.MockReadFolderDBWithErr)
			if !tc.MockReadFolderDBWithErr {
				setMockReadSubFolderDB(ts.mock, tc.MockReadSubFolderDBWithErr)
				if !tc.MockReadSubFolderDBWithErr {
					setMockFilesReadAllDB(ts.mock, tc.MockFilesReadAllDBWithErr)
				}
			}
		}

		// Act
		ts.handler.GetByID(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FolderTransactionSuite) TestReadFolderDB() {
	// Arrange
	setMockReadFolderDB(ts.mock, false)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FolderTransactionSuite) TestReadSubFolderDB() {
	// Arrange
	setMockReadSubFolderDB(ts.mock, false)

	// Act
	_, err := readSubFolderDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}
