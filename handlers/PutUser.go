package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	varsUpdate := mux.Vars(r)
	idStrUpdate := varsUpdate["id"]

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStrUpdate)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = db.QueryRow(
		context.Background(),
		"UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email",
		user.Name, user.Email, id,
	).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser) //Из Go в json
}
