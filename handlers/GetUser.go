package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetUserHandler возвращает список пользователей
// @Summary Получить список пользователей
// @Description Возвращает всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {string} string "Ошибка сервера"
// @Router /users [get]
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
