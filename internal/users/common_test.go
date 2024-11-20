package users

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTransactionSuite struct {
	suite.Suite
	conn    *sql.DB
	mock    sqlmock.Sqlmock
	handler handler
	entity  *User
}

func (ts *UserTransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.handler = handler{ts.conn}

	ts.entity = &User{
		Name:     "Teste user 1",
		Login:    "testuser1",
		Password: "testpassword",
	}
}

func (ts *UserTransactionSuite) AfterTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTransactionSuite))
}
