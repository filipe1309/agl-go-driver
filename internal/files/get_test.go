package files

import (
	"github.com/stretchr/testify/assert"
)

func (ts *FileTransactionSuite) TestReadDB() {
	// Arrange
	setMockReadDB(ts.mock)

	// Act
	_, err := ReadDB(ts.conn, 1)

	// Assert
	assert.NoError(ts.T(), err)
}
