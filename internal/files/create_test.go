package files

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	folder, err := New(1, "Test name", "webp", "path")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(0, 1, "Test name", "webp", "path", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
