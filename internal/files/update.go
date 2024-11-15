package files

import (
	"database/sql"
	"time"
)

func Update(db *sql.DB, id int64, file *File) (int64, error) {
	file.UpdatedAt = time.Now()

	stmt := `UPDATE files SET name = $1, updated_at = $2, deleted = $3 WHERE id = $4`
	result, err := db.Exec(stmt, file.Name, file.UpdatedAt, file.Deleted, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
