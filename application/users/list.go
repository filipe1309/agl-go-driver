package users

import (
	"encoding/json"
	"net/http"
)

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.factory.RestoreAll()
	if err != nil {
		// TODO: Check if the error is sql.ErrNoRows and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
