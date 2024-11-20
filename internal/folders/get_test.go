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

func (ts *FolderTransactionSuite) TestGetByID() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockReadFolderDB(ts.mock)
	setMockReadSubFolderDB(ts.mock)
	setMockFilesReadAllDB(ts.mock)

	// Act
	ts.handler.GetByID(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *FolderTransactionSuite) TestReadFolderDB() {
	// Arrange
	setMockReadFolderDB(ts.mock)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FolderTransactionSuite) TestReadSubFolderDB() {
	// Arrange
	setMockReadSubFolderDB(ts.mock)

	// Act
	_, err := readSubFolderDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadFolderDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, 2, "Test name", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)
}
