package folders

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestUpdate() {
	// Arrange
	folder := &Folder{
		ID:   1,
		Name: "Test folder 1",
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(folder)
	assert.NoError(ts.T(), err)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/folders/{id}", &b)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockUpdateDB(ts.mock, 1)
	setMockReadFolderDB(ts.mock)

	// Act
	ts.handler.Update(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *FolderTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1)

	// Act
	_, err := UpdateDB(ts.conn, 1, &Folder{
		Name: "Test folder 1",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET name = $1, updated_at = $2 WHERE id = $3`)).
		WithArgs("Test folder 1", sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
