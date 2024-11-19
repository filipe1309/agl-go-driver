package files

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestSoftDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db, nil, nil}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.SoftDelete(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

func TestSoftDeleteDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
		// WithArgs(AnyTime{}, 1).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = SoftDeleteDB(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
