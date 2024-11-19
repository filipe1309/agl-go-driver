package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

// type AnyTime struct{}
// func (a AnyTime) Match(v driver.Value) bool {
// 	_, ok := v.(time.Time)
// 	return ok
// }

func TestSoftDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	user := &User{
		ID:   1,
		Name: "Test user 1",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(user)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.SoftDelete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

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
