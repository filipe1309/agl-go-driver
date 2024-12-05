package users

import (
	"encoding/json"
	"net/http"

	domain "github.com/filipe1309/agl-go-driver/internal/users"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := domain.DecodeAndCreate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.InsertDB(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
