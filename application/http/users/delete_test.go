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
		IDStr                 string
		ID             int64
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
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockSoftDeleteDB(ts.mock, tc.ID, tc.MockWithErr)
		}

		// Act
		ts.handler.SoftDelete(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
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
