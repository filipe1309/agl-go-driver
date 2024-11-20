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
		ID                 string
		WithMock           bool
		MockUser           *User
		MockUpdatedWithErr bool
		ExpectedStatusCode int
	}{
		{"1", true, &User{ID: 1, Name: "Test user 1"}, false, http.StatusOK},
		{"A", false, nil, false, http.StatusBadRequest},
		{"25", true, &User{ID: 25, Name: "Test user 25"}, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		// Arrange

		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(tc.MockUser)
		assert.NoError(ts.T(), err)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/{id}", &b)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockUpdateDB(ts.mock, int(tc.MockUser.ID), tc.MockUpdatedWithErr)
			if !tc.MockUpdatedWithErr {
				setMockReadDB(ts.mock, int(tc.MockUser.ID))
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
