package folders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestSoftDelete() {
	tcs := []struct {
		Name                       string
		IDStr                      string
		WithMockDB                 bool
		MockWithFilesReadAllDBErr  bool
		MockWithFilesUpdateDBErr   bool
		MockWithFilesUpdateDBErr2  bool
		MockWithReadSubFolderDBErr bool
		MockWithSoftDeleteDBErr    bool
		ExpectedStatusCode         int
	}{
		{"Success", "1", true, false, false, false, false, false, http.StatusNoContent},
		{"Invalid url param - id", "A", false, false, false, false, false, false, http.StatusBadRequest},
		{"DB error - files read all", "25", true, true, false, false, false, false, http.StatusInternalServerError},
		{"DB error - files update 1", "25", true, false, true, false, false, false, http.StatusInternalServerError},
		{"DB error - files update 2", "25", true, false, false, true, false, false, http.StatusInternalServerError},
		{"DB error - read sub folder", "1", true, false, false, false, true, false, http.StatusInternalServerError},
		{"DB error - soft delete", "1", true, false, false, false, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/folders/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMockDB {
			setMockFilesReadAllDB(ts.mock, tc.MockWithFilesReadAllDBErr)
			if !tc.MockWithFilesReadAllDBErr {
				setMockFilesUpdateDB(ts.mock, "Test file 1", 1, true, tc.MockWithFilesUpdateDBErr)
				if !tc.MockWithFilesUpdateDBErr {
					setMockFilesUpdateDB(ts.mock, "Test file 2", 2, true, tc.MockWithFilesUpdateDBErr2)
					if !tc.MockWithFilesUpdateDBErr2 {
						setMockReadSubFolderDB(ts.mock, tc.MockWithReadSubFolderDBErr)
						if !tc.MockWithReadSubFolderDBErr {
							setSoftDeleteDB(ts.mock, tc.MockWithSoftDeleteDBErr)
						}
					}
				}
			}
		}

		// Act
		ts.handler.SoftDelete(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FolderTransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setSoftDeleteDB(ts.mock, false)

	// Act
	err := SoftDeleteDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockFilesUpdateDB(mock sqlmock.Sqlmock, fileName string, id int, deleted bool, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs(fileName, sqlmock.AnyArg(), deleted, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err {
		exp.WillReturnError(sqlmock.ErrCancelled)
	}
}

func setSoftDeleteDB(mock sqlmock.Sqlmock, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err {
		exp.WillReturnError(sqlmock.ErrCancelled)
	}
}
