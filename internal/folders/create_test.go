package folders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/filipe1309/agl-go-driver/internal/common"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestCreate() {
	// Arrange
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(ts.entity)
	assert.NoError(ts.T(), err)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/folders", &b)

	setMockInsertDB(ts.mock, ts.entity, nil)

	// Act
	ts.handler.Create(rr, req)

	// Assert
	assert.Equal(ts.T(), http.StatusCreated, rr.Code)
}

func (ts *FolderTransactionSuite) TestInsertRootDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity, nil)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FolderTransactionSuite) TestInsertWithFolderDB() {
	// Arrange
	ts.entity.ParentID = common.ValidNullInt64(1)

	setMockInsertDB(ts.mock, ts.entity, 1)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *Folder, parentID any) {
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`)).
		WithArgs(parentID, entity.Name, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
