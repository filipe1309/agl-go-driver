package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/filipe1309/agl-go-driver/internal/files"
	"github.com/go-chi/chi/v5"
)

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folder, err := ReadFolderDB(h.db, int64(id))
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := ReadFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folderContent := FolderContent{Folder: *folder, Content: content}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folderContent)
}

func ReadFolderDB(db *sql.DB, id int64) (*Folder, error) {
	stmt := `SELECT * FROM folders WHERE id = $1`

	var folder Folder
	row := db.QueryRow(stmt, id)
	err := row.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt, &folder.Deleted)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func readSubFolderDB(db *sql.DB, id int64) ([]Folder, error) {
	stmt := `SELECT * FROM folders WHERE parent_id = $1 AND deleted = FALSE`

	rows, err := db.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	folders := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.ParentID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt, &folder.Deleted)
		if err != nil {
			return nil, err
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func ReadFolderContent(db *sql.DB, folderID int64) ([]FolderResource, error) {
	subFolders, err := readSubFolderDB(db, folderID)
	if err != nil {
		return nil, err
	}

	resources := make([]FolderResource, 0)
	for _, folder := range subFolders {
		resources = append(resources, FolderResource{
			ID:        folder.ID,
			Name:      folder.Name,
			Type:      "directory",
			CreatedAt: folder.CreatedAt,
			UpdatedAt: folder.UpdatedAt,
		})
	}

	folderFiles, err := files.ReadAllDB(db, folderID)
	if err != nil {
		return nil, err
	}

	for _, file := range folderFiles {
		resources = append(resources, FolderResource{
			ID:        file.ID,
			Name:      file.Name,
			Type:      file.Type,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
		})
	}

	return resources, nil
}
