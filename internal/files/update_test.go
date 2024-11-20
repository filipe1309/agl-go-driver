package files

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

func (ts *FileTransactionSuite) TestUpdate() {
	// Arrange
	file := &File{
		ID:   1,
		Name: "Test file 2",
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(file)
	assert.NoError(ts.T(), err)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/files/{id}", &b)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockReadDB(ts.mock)
	setMockUpdateDB(ts.mock, 1)

	// Act
	ts.handler.Update(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusOK, rr.Code)
}

func (ts *FileTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1)

	// Act
	_, err := UpdateDB(ts.conn, 1, &File{
		Name: "Test file 2",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int) {
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`)).
		WithArgs("Test file 2", sqlmock.AnyArg(), false, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
