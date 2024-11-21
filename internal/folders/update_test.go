package folders

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *FolderTransactionSuite) TestUpdate() {
	tcs := []struct {
		Name                    string
		ID                      int
		IDStr                   string
		WithMock                bool
		MockFolder              *Folder
		MockUpdatedDBWithErr    bool
		MockReadFolderDBWithErr bool
		ExpectedStatusCode      int
	}{
		{
			Name:                    "Success",
			ID:                      1,
			IDStr:                   "1",
			WithMock:                true,
			MockFolder:              &Folder{ID: 1, Name: "Test folder 1"},
			MockUpdatedDBWithErr:    false,
			MockReadFolderDBWithErr: false,
			ExpectedStatusCode:      http.StatusOK,
		},
		{
			Name:                    "Invalid folder - no name",
			ID:                      2,
			IDStr:                   "2",
			WithMock:                false,
			MockFolder:              &Folder{ID: 2},
			MockUpdatedDBWithErr:    false,
			MockReadFolderDBWithErr: false,
			ExpectedStatusCode:      http.StatusBadRequest,
		},
		{
			Name:                    "Invalid json - empty body",
			ID:                      3,
			IDStr:                   "3",
			WithMock:                false,
			MockFolder:              &Folder{},
			MockUpdatedDBWithErr:    false,
			MockReadFolderDBWithErr: false,
			ExpectedStatusCode:      http.StatusBadRequest,
		},
		{
			Name:                    "Invalid url param - id",
			ID:                      -1,
			IDStr:                   "A",
			WithMock:                false,
			MockFolder:              &Folder{ID: 1, Name: "Test folder 1"},
			MockUpdatedDBWithErr:    false,
			MockReadFolderDBWithErr: false,
			ExpectedStatusCode:      http.StatusBadRequest,
		},
		{
			Name:                    "DB error - update",
			ID:                      1,
			IDStr:                   "1",
			WithMock:                true,
			MockFolder:              &Folder{ID: 1, Name: "Test folder 1"},
			MockUpdatedDBWithErr:    true,
			MockReadFolderDBWithErr: false,
			ExpectedStatusCode:      http.StatusInternalServerError,
		},

		{
			Name:                    "DB error - no read id",
			ID:                      1,
			IDStr:                   "1",
			WithMock:                true,
			MockFolder:              &Folder{ID: 1, Name: "Test folder 1"},
			MockUpdatedDBWithErr:    false,
			MockReadFolderDBWithErr: true,
			ExpectedStatusCode:      http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		// Arrange
		var b bytes.Buffer

		if tc.Name != "Invalid json - empty body" {
			err := json.NewEncoder(&b).Encode(tc.MockFolder)
			assert.NoError(ts.T(), err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/folders/{id}", &b)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.IDStr)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		fmt.Println(tc.Name)
		if tc.WithMock {
			fmt.Println("WithMock", tc.MockUpdatedDBWithErr)
			setMockUpdateDB(ts.mock, tc.ID, tc.MockUpdatedDBWithErr)
			if !tc.MockUpdatedDBWithErr {
				setMockReadFolderDB(ts.mock, tc.MockReadFolderDBWithErr)
			}
		}

		// Act
		ts.handler.Update(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *FolderTransactionSuite) TestUpdateDB() {
	// Arrange
	setMockUpdateDB(ts.mock, 1, false)

	// Act
	_, err := UpdateDB(ts.conn, 1, &Folder{
		Name: "Test folder 1",
	})

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockUpdateDB(mock sqlmock.Sqlmock, id int, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET name = $1, updated_at = $2 WHERE id = $3`)).
		WithArgs("Test folder 1", sqlmock.AnyArg(), id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
