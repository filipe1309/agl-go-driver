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
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockFilesReadAllDB(ts.mock)
	setMockFilesUpdateDB(ts.mock, "Test file 1", 1, true)
	setMockFilesUpdateDB(ts.mock, "Test file 2", 2, true)
	setMockReadSubFolderDB(ts.mock)
	setSoftDeleteDB(ts.mock)

	// Act
	ts.handler.SoftDelete(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusNoContent, rr.Code)
}

func (ts *FolderTransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setSoftDeleteDB(ts.mock)

	// Act
	err := SoftDeleteDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockFilesUpdateDB(mock sqlmock.Sqlmock, fileName string, id int, deleted bool) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs(fileName, sqlmock.AnyArg(), deleted, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func setSoftDeleteDB(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
