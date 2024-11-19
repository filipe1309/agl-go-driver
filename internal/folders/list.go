package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/filipe1309/agl-go-driver/internal/files"
)

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	content, err := ReadRootFolderContent(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	folderContent := FolderContent{Folder: Folder{Name: "root"}, Content: content}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folderContent)
}

func readRootSubFolderDB(db *sql.DB) ([]Folder, error) {
	stmt := `SELECT * FROM folders WHERE parent_id IS NULL AND deleted = FALSE`

	rows, err := db.Query(stmt)
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

func ReadRootFolderContent(db *sql.DB) ([]FolderResource, error) {
	subFolders, err := readRootSubFolderDB(db)
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

	folderFiles, err := files.ReadAllRootDB(db)
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
