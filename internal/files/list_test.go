package files

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

	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(rows)

	_, err = ReadAll(db, 0)
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

func TestReadAllRoot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)

	_, err = ReadAllRoot(db)
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
