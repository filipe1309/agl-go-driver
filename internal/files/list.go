package files

import (
	"database/sql"
	"log"
)

func ReadAllDB(db *sql.DB, folder_id int64) ([]File, error) {
	stmt := `SELECT * FROM files WHERE folder_id = $1 AND deleted = FALSE`
	rows, err := db.Query(stmt, folder_id)
	if err != nil {
		return nil, err
	}
	return selectAllFiles(rows)
}

func ReadAllRootDB(db *sql.DB) ([]File, error) {
	stmt := `SELECT * FROM files WHERE folder_id IS NULL AND deleted = FALSE`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	return selectAllFiles(rows)
}

func selectAllFiles(rows *sql.Rows) ([]File, error) {
	defer rows.Close()

	files := make([]File, 0)
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.FolderID, &file.OwnerID, &file.Name, &file.Type, &file.Path, &file.CreatedAt, &file.UpdatedAt, &file.Deleted)
		if err != nil {
			log.Println(err)
			continue
		}

		files = append(files, file)
	}

	return files, nil
}
