package users

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domain "github.com/filipe1309/agl-go-driver/internal/users"
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
		{
			Name:                "Invalid json",
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
		if tc.Name != "Invalid json" {
			err := json.NewEncoder(&b).Encode(ts.entity)
			assert.NoError(ts.T(), err)
		}

		ts.entity.ChangePassword(ts.entity.Password)

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



func setMockInsertDB(mock sqlmock.Sqlmock, entity *domain.User, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (name, login, password, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`)).
		WithArgs(entity.Name, entity.Login, sqlmock.AnyArg(), sqlmock.AnyArg())

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	}
}
