package users

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	user, err := New("Test name", "testlogin", "testpassword")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`INSERT INTO users (name, login, password, updated_at)*`).
		WithArgs("Test name", "testlogin", user.Password, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, user)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
