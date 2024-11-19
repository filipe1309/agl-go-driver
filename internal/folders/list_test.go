package folders

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, "Test folder 1", time.Now(), time.Now(), false).
		AddRow(5, nil, "Test folder 1", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "updated_at", "deleted"}).
		AddRow(1, 1, 1, "Test name", "testtype", "testpath", time.Now(), time.Now(), false).
		AddRow(2, 1, 1, "Test name 2", "testtype", "testpath", time.Now(), time.Now(), true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows2)

	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		fmt.Println(rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestReadRootSubFolderDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "updated_at", "deleted"}).
		AddRow(1, nil, "Test folder 1", time.Now(), time.Now(), false).
		AddRow(5, nil, "Test folder 1", time.Now(), time.Now(), false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM folders WHERE parent_id IS NULL AND deleted = FALSE`)).
		WillReturnRows(rows)

	_, err = readRootSubFolderDB(db)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
