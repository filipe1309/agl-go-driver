package folders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	folder, err := New("Test name", 0)
	if err != nil {
		t.Error(err)
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(folder)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/folders", &b)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`)).
		WithArgs(nil, "Test name", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

func TestInsertRoot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	folder, err := New("Test name", 0)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`)).
		WithArgs(nil, "Test name", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestInsertWithFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	folder, err := New("Test name", 1)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`)).
		WithArgs(1, "Test name", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
