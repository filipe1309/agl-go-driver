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
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestUpdate() {
	tcs := []struct {
		Name               string
		ID                 int
		IDStr              string
		WithMock           bool
		MockUser           *User
		MockUpdatedWithErr bool
		MockReadWithErr    bool
		ExpectedStatusCode int
	}{
		{
			Name:               "Success",
			ID:                 1,
			IDStr:              "1",
			WithMock:           true,
			MockUser:           &User{ID: 1, Name: "Test user 1"},
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "Invalid user - no name",
			ID:                 2,
			IDStr:              "2",
			WithMock:           false,
			MockUser:           &User{ID: 2},
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid json - empty body",
			ID:                 3,
			IDStr:              "3",
			WithMock:           false,
			MockUser:           &User{},
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "Invalid url param - id",
			ID:                 -1,
			IDStr:              "A",
			WithMock:           false,
			MockUser:           &User{Name: "Test user A"},
			MockUpdatedWithErr: false,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "DB error - no Update id",
			ID:                 25,
			IDStr:              "25",
			WithMock:           true,
			MockUser:           &User{ID: 25, Name: "Test user 25"},
			MockUpdatedWithErr: true,
			MockReadWithErr:    false,
			ExpectedStatusCode: http.StatusInternalServerError,
		},

		{
			Name:               "DB error - no read id",
			ID:                 26,
			IDStr:              "26",
			WithMock:           true,
			MockUser:           &User{ID: 26, Name: "Test user 26"},
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

		if tc.WithMock {
			setMockUpdateDB(ts.mock, tc.ID, tc.MockUpdatedWithErr)
			if !tc.MockUpdatedWithErr {
				setMockReadDB(ts.mock, tc.ID, tc.MockReadWithErr)
			}
		}

		// Act
		ts.handler.Update(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *UserTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1, false)

	// Act
	_, err := UpdateDB(ts.conn, 1, &User{
		Name: "Test user 1",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, updated_at = $2 WHERE id = $3`)).
		WithArgs(fmt.Sprintf("%s %d", "Test user", id), sqlmock.AnyArg(), id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
