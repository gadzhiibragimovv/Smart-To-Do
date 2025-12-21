package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	varsDelete := mux.Vars(r)
	idStrDelete := varsDelete["id"]

	id, err := strconv.Atoi(idStrDelete)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed from delete user:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}
