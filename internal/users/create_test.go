package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	// Arrange
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(ts.entity)
	assert.NoError(ts.T(), err)

	ts.entity.SetPassword(ts.entity.Password)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", &b)

	setMockInsertDB(ts.mock, ts.entity)

	// Act
	ts.handler.Create(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusCreated, rr.Code)
}

func (ts *TransactionSuite) TestInsertDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *User) {
	mock.ExpectExec(`INSERT INTO users (name, login, password, updated_at)*`).
		WithArgs(entity.Name, entity.Login, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
