package folders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

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

func setMockFilesReadAllDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test file 1", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test file 2", "testtype", "testpath", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(rows)
}

func setMockFilesUpdateDB(mock sqlmock.Sqlmock, fileName string, id int, deleted bool) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs(fileName, sqlmock.AnyArg(), deleted, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func setMockReadSubFolderDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(2, 3, "Test sub name", time.Now(), time.Now(), false).
		AddRow(3, 3, "Test sub name 2", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`)).
		WithArgs(1).
		WillReturnRows(rows)
}

func setSoftDeleteDB(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
