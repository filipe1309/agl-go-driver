package files

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

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db, nil, nil}

	file := &File{
		ID:   1,
		Name: "Test file 2",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(file)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/files/{id}", &b)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// ReadDB
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 3, "Test file 1", "testtype", "testpath", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	// UpdateDB
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs("Test file 2", sqlmock.AnyArg(), false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Update(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

func TestUpdateDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs("Test name", sqlmock.AnyArg(), false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = UpdateDB(db, 1, &File{
		Name: "Test name",
	})
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
