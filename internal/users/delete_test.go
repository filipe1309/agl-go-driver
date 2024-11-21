package users

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestSoftDelete() {
	tcs := []struct {
		Name 						 string
		ID                 string
		MockID             int64
		WithMock           bool
		MockWithErr        bool
		ExpectedStatusCode int
	}{
		{"Success", "1", 1, true, false, http.StatusNoContent},
		{"Invalid url param - id", "A", -1, false, true, http.StatusBadRequest},
		{"DB error", "25", 25, true, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/users/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockSoftDeleteDB(ts.mock, tc.MockID, tc.MockWithErr)
		}

		// Act
		ts.handler.SoftDelete(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *UserTransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setMockSoftDeleteDB(ts.mock, 1, false)

	// Act
	err := SoftDeleteDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockSoftDeleteDB(mock sqlmock.Sqlmock, id int64, err bool) {

	exp := mock.ExpectExec(`UPDATE users SET *`).
		// WithArgs(AnyTime{}, id).
		WithArgs(sqlmock.AnyArg(), id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
