package files

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestSoftDelete() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockSoftDeleteDB(ts.mock)

	// Act
	ts.handler.SoftDelete(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusNoContent, rr.Code)
}

func (ts *FileTransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setMockSoftDeleteDB(ts.mock)

	// Act
	err := SoftDeleteDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockSoftDeleteDB(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
