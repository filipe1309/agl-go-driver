package users

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestList() {
	// Arrange
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	setMockReadAllDB(ts.mock)

	// Act
	ts.handler.List(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *UserTransactionSuite) TestReadAllDB() {
	// Arrange
	setMockReadAllDB(ts.mock)

	// Act
	_, err := ReadAllDB(ts.conn)

	// Assert
	assert.NoError(ts.T(), err)

	// if len(list) > 1 { // this doesn't work with sqlmock, because it doesn't filter the rows
	// 	t.Error("Expected 1 row, got", len(list))
	// }
}

func setMockReadAllDB(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test name", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false).
		AddRow(2, "Test name 2", "testlogin2", "testpassword", time.Now(), time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE deleted = FALSE`)).
		WillReturnRows(rows)
}
