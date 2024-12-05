package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/common"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	domain "github.com/filipe1309/agl-go-driver/internal/users"
)

func (ts *UserTransactionSuite) TestUpdate() {
	tcs := []struct {
		Name               string
		IDStr              string
		MockUser           *domain.User
		WithReadDBMock     bool
		WithUpdateDBMock   bool
		MockUpdatedWithErr bool
		MockReadWithErr    bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Success",
			IDStr:              "1",
			MockUser:           &domain.User{ID: 1, Name: "Test user 1"},
			WithReadDBMock:     true,
			WithUpdateDBMock:   true,
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "Invalid user - no name",
			IDStr:              "2",
			WithReadDBMock:     true,
			WithUpdateDBMock:   false,
			MockUser:           &domain.User{ID: 2},
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid json - empty body",
			IDStr:              "0",
			MockUser:           &domain.User{},
			WithReadDBMock:     true,
			WithUpdateDBMock:   false,
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid url param - id",
			IDStr:              "A",
			MockUser:           &domain.User{ID: -1, Name: "Test user A"},
			WithReadDBMock:     false,
			WithUpdateDBMock:   false,
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "DB error - no Update id",
			IDStr:              "1",
			MockUser:           &domain.User{ID: 1, Name: "Test user 1"},
			WithReadDBMock:     true,
			WithUpdateDBMock:   true,
			MockUpdatedWithErr: true,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "DB error - no read id",
			IDStr:              "26",
			WithReadDBMock:     true,
			WithUpdateDBMock:   false,
			MockUser:           &domain.User{ID: 26, Name: "Test user 26"},
			MockUpdatedWithErr: false,
			MockReadWithErr:    true,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		var b bytes.Buffer
		if tc.Name != "Invalid json - empty body" {
			err := json.NewEncoder(&b).Encode(tc.MockUser)
			assert.NoError(ts.T(), err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/{id}", &b)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(common.SetUserIDInContext(req.Context(), tc.MockUser.ID))

		if tc.WithReadDBMock {
			setMockReadDB(ts.mock, int(tc.MockUser.ID), tc.MockReadWithErr)
			if tc.WithUpdateDBMock && !tc.MockReadWithErr {
				setMockUpdateDB(ts.mock, int(tc.MockUser.ID), tc.MockUpdatedWithErr)
			}
		}

		// Act
		ts.handler.Update(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, updated_at = $2, last_login = $3 WHERE id = $4`)).
		WithArgs(fmt.Sprintf("%s %d", "Test user", id), sqlmock.AnyArg(), sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}
