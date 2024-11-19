package folders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestSoftDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	folder := &Folder{
		ID:   1,
		Name: "Test folder 1",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(folder)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(2, 3, "Test sub name", time.Now(), time.Now(), false).
		AddRow(3, 3, "Test sub name 2", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`)).
		WithArgs(1).
		WillReturnRows(rows2)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
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

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET updated_at = $1, deleted = TRUE WHERE id = $2`)).
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
