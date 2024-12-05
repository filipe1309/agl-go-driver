package users

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/factories"
	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/filipe1309/agl-go-driver/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTransactionSuite struct {
	suite.Suite
	conn    *sql.DB
	mock    sqlmock.Sqlmock
	handler handler
	entity  *users.User
}

func (ts *UserTransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	repo := repositories.NewUserRepository(ts.conn)
	factory := factories.NewUserFactory(repo)
	ts.handler = handler{repo, factory}

	ts.entity = &users.User{
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

func setMockReadDB(mock sqlmock.Sqlmock, id int, err bool) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test user 1", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false)
	exp := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE id = $1`)).
		WithArgs(id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnRows(rows)
	}
}
