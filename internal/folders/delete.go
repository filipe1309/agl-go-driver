package folders

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/filipe1309/agl-go-driver/internal/files"
	"github.com/go-chi/chi/v5"
)

func (h *handler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: List folders

	err = deleteFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = SoftDelete(h.db, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	return deleteSubFolders(db, folderID)
}

func deleteSubFolders(db *sql.DB, folderID int64) error {
	subFoldersList, err := readSubFolder(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subFoldersList))
	for _, subFolder := range subFoldersList {
		subFolder.Deleted = true
		err = SoftDelete(db, subFolder.ID)
		if err != nil {
			break
		}
		err = deleteFolderContent(db, subFolder.ID)
		if err != nil {
			// subFolder.Deleted = false
			Update(db, subFolder.ID, &subFolder)
			break
		}

		removedFolders = append(removedFolders, subFolder)
	}

	if len(subFoldersList) != len(removedFolders) {
		for _, subFolder := range removedFolders {
			// subFolder.Deleted = false
			_, _ = Update(db, subFolder.ID, &subFolder)
		}
	}

	return nil
}

func deleteFiles(db *sql.DB, folderID int64) error {
	filesList, err := files.ReadAll(db, folderID)
	if err != nil {
		return err
	}

	removedFiles := make([]files.File, len(filesList))
	for _, file := range filesList {
		file.Deleted = true
		_, err = files.Update(db, file.ID, &file)
		if err != nil {
			break
		}
		removedFiles = append(removedFiles, file)
	}

	if len(filesList) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			_, _ = files.Update(db, file.ID, &file)
		}

		return err
	}

	return nil
}

func SoftDelete(db *sql.DB, id int64) error {
	stmt := `UPDATE folders SET updated_at = $1, deleted = true WHERE id = $2`

	_, err := db.Exec(stmt, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
