package users

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestGetByID() {
	tcs := []struct {
		Name               string
		ID                 int
		IDStr              string
		WithMockDB         bool
		MockReadDBWithErr  bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Success",
			ID:                 1,
			IDStr:              "1",
			WithMockDB:         true,
			MockReadDBWithErr:  false,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "DB error",
			ID:                 1,
			IDStr:              "1",
			WithMockDB:         true,
			MockReadDBWithErr:  true,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "Invalid url param - id",
			ID:                 0,
			IDStr:              "A",
			WithMockDB:         false,
			MockReadDBWithErr:  false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMockDB {
			setMockReadDB(ts.mock, tc.ID, tc.MockReadDBWithErr)
		}

		// Act
		ts.handler.GetByID(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}
