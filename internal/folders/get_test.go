package folders

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/folders/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	// ReadFolderDB
	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, 2, "Test name", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	// readSubFolderDB
	rows2 := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(2, 3, "Test sub name", time.Now(), time.Now(), false).
		AddRow(3, 3, "Test sub name 2", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`)).
		WithArgs(1).
		WillReturnRows(rows2)

	// files.ReadAllDB
	fileRows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`)).
		WillReturnRows(fileRows)

	h.GetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

func TestReadFolderDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, 2, "Test name", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = ReadDB(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestReadSubFolderDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(2, 3, "Test sub name", time.Now(), time.Now(), false).
		AddRow(3, 3, "Test sub name 2", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = readSubFolderDB(db, 1)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
