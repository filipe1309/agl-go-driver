package users

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test name", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = Read(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
