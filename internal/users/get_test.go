package users

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestGetByID() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockReadDB(ts.mock, 1)

	// Act
	ts.handler.GetByID(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *UserTransactionSuite) TestReadDB() {
	// Arrange
	setMockReadDB(ts.mock, 1)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}
