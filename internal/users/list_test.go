package users

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReadAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "updated_at", "last_login", "deleted"}).
		AddRow(1, "Test name", "testlogin", "testpassword", time.Now(), time.Now(), time.Now(), false).
		AddRow(2, "Test name 2", "testlogin2", "testpassword", time.Now(), time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE deleted = FALSE`)).
		WillReturnRows(rows)

	_, err = ReadAll(db)
	if err != nil {
		t.Error(err)
	}

	// if len(list) > 1 { // this doesn't work with sqlmock, because it doesn't filter the rows
	// 	t.Error("Expected 1 row, got", len(list))
	// }

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
