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
	tcs := []struct {
		Name                string
		FolderID            any
		WithMockDB          bool
		MockInsertDBWithErr bool
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success - insert on root",
			FolderID:            nil,
			WithMockDB:          true,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusCreated,
		},
		{
			Name:                "Success - insert on existent folder",
			FolderID:            int64(1),
			WithMockDB:          true,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusCreated,
		},
		{
			Name:                "DB error",
			FolderID:            int64(1),
			WithMockDB:          true,
			MockInsertDBWithErr: true,
			ExpectedStatusCode:  http.StatusInternalServerError,
		},
		{
			Name:                "Invalid user - empty name",
			FolderID:            int64(1),
			WithMockDB:          false,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
		{
			Name:                "Invalid json",
			FolderID:            int64(1),
			WithMockDB:          false,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		// Arrange
		if tc.FolderID != nil {
			ts.entity.ParentID = common.ValidNullInt64(tc.FolderID.(int64))
		}

		if tc.Name == "Invalid user - empty name" {
			ts.entity.Name = ""
		}
		var b bytes.Buffer
		if tc.Name != "Invalid json" {
			err := json.NewEncoder(&b).Encode(ts.entity)
			assert.NoError(ts.T(), err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/folders", &b)

		if tc.WithMockDB {
			setMockInsertDB(ts.mock, ts.entity, tc.FolderID, tc.MockInsertDBWithErr)
		}

		// Act
		ts.handler.Create(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FolderTransactionSuite) TestInsertRootDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity, nil, false)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func (ts *FolderTransactionSuite) TestInsertWithFolderDB() {
	// Arrange
	ts.entity.ParentID = common.ValidNullInt64(1)

	setMockInsertDB(ts.mock, ts.entity, 1, false)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *Folder, parentID any, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders (parent_id, name, updated_at) VALUES ($1, $2, $3)`)).
		WithArgs(parentID, entity.Name, sqlmock.AnyArg())

	if err {
		exp.WillReturnError(sqlmock.ErrCancelled)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
