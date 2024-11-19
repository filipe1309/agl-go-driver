package files

import "database/sql"

func ReadDB(db *sql.DB, id int64) (*File, error) {
	stmt := `SELECT * FROM files WHERE id = $1`

	row := db.QueryRow(stmt, id)

	var file File
	err := row.Scan(&file.ID, &file.FolderID, &file.OwnerID, &file.Name, &file.Type, &file.Path, &file.CreatedAt, &file.UpdatedAt, &file.Deleted)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
