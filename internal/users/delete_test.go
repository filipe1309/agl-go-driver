package users

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// type AnyTime struct{}
// func (a AnyTime) Match(v driver.Value) bool {
// 	_, ok := v.(time.Time)
// 	return ok
// }

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET *`).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = SoftDelete(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
