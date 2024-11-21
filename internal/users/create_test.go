package users

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *UserTransactionSuite) TestCreate() {
	tcs := []struct {
		Name                string
		Password            string
		WithMockDB          bool
		MockInsertDBWithErr bool
		ExpectedStatusCode  int
	}{
		{
			Name:                "Success",
			Password:            ts.entity.Password,
			WithMockDB:          true,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusCreated,
		},
		{
			Name:                "DB error",
			Password:            ts.entity.Password,
			WithMockDB:          true,
			MockInsertDBWithErr: true,
			ExpectedStatusCode:  http.StatusInternalServerError,
		},
		{
			Name:                "Invalid user - wrong password",
			Password:            "",
			WithMockDB:          false,
			MockInsertDBWithErr: false,
			ExpectedStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		// Arrange
		ts.entity.Password = tc.Password
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(ts.entity)
		assert.NoError(ts.T(), err)

		ts.entity.SetPassword(ts.entity.Password)
		
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", &b)

		if tc.WithMockDB {
			setMockInsertDB(ts.mock, ts.entity, tc.MockInsertDBWithErr)
		}

		// Act
		ts.handler.Create(rr, req)

		// Assert
		assert.Equal(ts.T(), tc.ExpectedStatusCode, rr.Code)
	}
}

func (ts *UserTransactionSuite) TestInsertDB() {
	// Arrange
	setMockInsertDB(ts.mock, ts.entity, false)

	// Act
	_, err := InsertDB(ts.conn, ts.entity)

	// Assert
	assert.NoError(ts.T(), err)
}

func setMockInsertDB(mock sqlmock.Sqlmock, entity *User, err bool) {
	exp := mock.ExpectExec(`INSERT INTO users (name, login, password, updated_at)*`).
		WithArgs(entity.Name, entity.Login, sqlmock.AnyArg(), sqlmock.AnyArg())

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
