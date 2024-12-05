package users

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestList() {
	tcs := []struct {
		Name                 string
		MockReadAllDBWithErr bool
		ExpectedStatusCode   int
	}{
		{
			Name:                 "Success",
			MockReadAllDBWithErr: false,
			ExpectedStatusCode:   http.StatusOK,
		},
		{
			Name:                 "DB error",
			MockReadAllDBWithErr: true,
			ExpectedStatusCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)

		setMockReadAllDB(ts.mock, tc.MockReadAllDBWithErr)

		// Act
		ts.handler.List(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func setMockReadAllDB(mock sqlmock.Sqlmock, err bool) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test name", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false).
		AddRow(2, "Test name 2", "testlogin2", "testpassword", time.Now(), time.Now(), time.Now(), true)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE deleted = FALSE`))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnRows(rows)
	}
}
