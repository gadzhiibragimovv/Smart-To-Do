package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{}

	rows, err := db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Failed to query users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			http.Error(w, "Error scanning user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "error missing in path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow(
		context.Background(),
		"SELECT id, name FROM users WHERE id = $1",
		id,
	).Scan(&user.Id, &user.Name)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
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
