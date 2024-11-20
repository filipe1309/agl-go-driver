package users

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

func (ts *TransactionSuite) TestUpdate() {
	// Arrange
	user := &User{
		ID:   1,
		Name: "Test user 1",
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(user)
	assert.NoError(ts.T(), err)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/users/{id}", &b)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockUpdateDB(ts.mock, 1)
	setMockReadDB(ts.mock)

	// Act
	ts.handler.Update(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *TransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1)

	// Act
	_, err := UpdateDB(ts.conn, 1, &User{
		Name: "Test user 1",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, updated_at = $2 WHERE id = $3`)).
		WithArgs("Test user 1", sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
