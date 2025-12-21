package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		context.Background(),
		"INSERT INTO users (name, email) VALUES ($1,$2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.Id)

	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
