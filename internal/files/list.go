package files

import (
	"database/sql"
	"log"
)

func ReadAll(db *sql.DB, folder_id int64) ([]File, error) {
	stmt := `SELECT * FROM files WHERE folder_id = $1 AND deleted = false`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]File, 0)
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.ID, &file.FolderID, &file.OwnerID, &file.Name, &file.Type, &file.Path, &file.CreatedAt, &file.UpdatedAt, &file.Deleted)
		if err != nil {
			log.Println(err)
			continue
		}

		files = append(files, file)
	}

	return files, nil
}
