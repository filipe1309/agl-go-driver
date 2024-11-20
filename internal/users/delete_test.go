package users

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestSoftDelete() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockSoftDeleteDB(ts.mock, 1)

	// Act
	ts.handler.SoftDelete(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusNoContent, rr.Code)
}

func (ts *TransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setMockSoftDeleteDB(ts.mock, 1)

	// Act
	err := SoftDeleteDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockSoftDeleteDB(mock sqlmock.Sqlmock, id int) {
	mock.ExpectExec(`UPDATE users SET *`).
		// WithArgs(AnyTime{}, id).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
