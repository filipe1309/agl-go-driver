package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	user, err := New("Teste user 1", "testuser1", "testpassword")
	if err != nil {
		t.Error(err)
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(user)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", &b)
	mock.ExpectExec(`INSERT INTO users (name, login, password, updated_at)*`).
		WithArgs(user.Name, user.Login, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

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
