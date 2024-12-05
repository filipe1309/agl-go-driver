package repositories

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTransactionSuite struct {
	suite.Suite
	conn   *sql.DB
	mock   sqlmock.Sqlmock
	repo   UserRepository
	entity *users.User
}

func (ts *UserTransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.repo = UserRepository{ts.conn}

	ts.entity = &users.User{
		Name:     "Teste user 1",
		Login:    "testuser1",
		Password: "testpassword",
	}
}

func (ts *UserTransactionSuite) AfterUserTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTransactionSuite))
}

func (ts *UserTransactionSuite) TestInsertDB() {
	// Arrange
	setMockUserInsertDB(ts.mock, ts.entity, false)

	// Act
	_, err := ts.repo.InsertDB(ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

// func (ts *UserTransactionSuite) TestReadDB() {
// 	// Arrange
// 	setMockUserReadDB(ts.mock, 1, false)

// 	// Act
// 	_, err := ts.repo.ReadDB(1)

// 	// Assert
// 	assert.NoError(ts.T(), err)
// }

func (ts *UserTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUserUpdateDB(ts.mock, 1, false)

	// Act
	_, err := ts.repo.UpdateDB(1, &users.User{
		Name: "Test user 1",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *UserTransactionSuite) TestSoftDeleteDB() {
	// Arrange
	setMockUserSoftDeleteDB(ts.mock, 1, false)

	// Act
	err := ts.repo.SoftDeleteDB(1)

	// Assert
	assert.NoError(ts.T(), err)
}

// func (ts *UserTransactionSuite) TestReadAllDB() {
// 	// Arrange
// 	setMockUserReadAllDB(ts.mock, false)

// 	// Act
// 	row := ts.repo.ReadAllDB()
// 	ts.factory.RestoreAll()

// 	// Assert
// 	assert.NoError(ts.T(), err)

// 	// if len(list) > 1 { // this doesn't work with sqlmock, because it doesn't filter the rows
// 	// 	t.Error("Expected 1 row, got", len(list))
// 	// }
// }

func setMockUserInsertDB(mock sqlmock.Sqlmock, entity *users.User, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`)).
		WithArgs(entity.Name, entity.Login, sqlmock.AnyArg(), sqlmock.AnyArg())

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	}
}

func setMockUserReadDB(mock sqlmock.Sqlmock, id int, err bool) {
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

func setMockUserUpdateDB(mock sqlmock.Sqlmock, id int, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET name = $1, updated_at = $2, last_login = $3 WHERE id = $4`)).
		WithArgs(fmt.Sprintf("%s %d", "Test user", id), sqlmock.AnyArg(), sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	}
}

func setMockUserSoftDeleteDB(mock sqlmock.Sqlmock, id int64, err bool) {

	exp := mock.ExpectExec(`UPDATE users SET *`).
		// WithArgs(AnyTime{}, id).
		WithArgs(sqlmock.AnyArg(), id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func setMockUserReadAllDB(mock sqlmock.Sqlmock, err bool) {
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
