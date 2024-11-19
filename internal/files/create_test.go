package files

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/common"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db, nil, nil}

	b := new(bytes.Buffer)

	mw := multipart.NewWriter(b)

	file, err := os.Open("./testdata/test-image-1.jpg")
	if err != nil {
		t.Error(err)
	}

	fw, err := mw.CreateFormFile("file", "test-image-1.jpg")
	if err != nil {
		t.Error(err)
	}

	if _, err := io.Copy(fw, file); err != nil {
		t.Error(err)
	}

	mw.Close()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/files", b)
	// req.Header.Set("Content-Type", mw.FormDataContentType())
	// req.Header.Set("Content-Type", "multipart/form-data; boundary="+mw.Boundary())

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(nil, 1, "test-image-1.jpg", "image/jpg", "/", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		fmt.Println(rr.Body.String())
	}
}

func TestInsertRootDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	file, err := New(1, "test-image-1.jpg", "image/jpg", "/")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(nil, 1, "test-image-1.jpg", "image/jpg", "/", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = InsertDB(db, file)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestInsertWithFolderDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	file, err := New(1, "test-image-1.jpg", "image/jpg", "/")
	if err != nil {
		t.Error(err)
	}
	file.FolderID = common.ValidNullInt64(2)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(int64(2), 1, "test-image-1.jpg", "image/jpg", "/", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = InsertDB(db, file)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
