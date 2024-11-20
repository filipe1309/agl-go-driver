package users

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

func (ts *TransactionSuite) TestGetByID() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockReadDB(ts.mock)

	// Act
	ts.handler.GetByID(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *TransactionSuite) TestReadDB() {
	// Arrange
	setMockReadDB(ts.mock)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockReadDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test name", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)
}
