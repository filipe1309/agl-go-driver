package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/filipe1309/agl-go-driver/internal/common"
	"github.com/filipe1309/agl-go-driver/internal/queue"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // 32MB

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := fmt.Sprintf("/%s", fileHeader.Filename)

	err = h.bucket.Upload(file, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: replace hardcoded ownerID with session user
	ownerID := 1
	fileEntity, err := New(int64(ownerID), fileHeader.Filename, fileHeader.Header.Get("Content-Type"), path)
	if err != nil {
		h.bucket.Delete(path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folderID := r.Form.Get("folder_id")
	if folderID != "" {
		folderIDInt, err := strconv.Atoi(folderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fileEntity.FolderID = common.NullInt64{Int64: int64(folderIDInt), Valid: true}
	}

	id, err := InsertDB(h.db, fileEntity)
	if err != nil {
		// h.bucket.Delete(path)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileEntity.ID = id

	dto := queue.QueueDTO{
		Filename: fileHeader.Filename,
		Path:     path,
		ID:       int(id),
	}

	msg, err := dto.Marshal()
	if err != nil {
		// TODO: rollback
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.queue.Publish(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileEntity)
}

func InsertDB(db *sql.DB, file *File) (int64, error) {
	stmt := `INSERT INTO files (folder_id, owner_id, name, type, path, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	var folderID any
	if file.FolderID.Valid {
		folderID = file.FolderID.Int64
	} else {
		folderID = nil
	}

	result, err := db.Exec(stmt, folderID, file.OwnerID, file.Name, file.Type, file.Path, file.UpdatedAt)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}
