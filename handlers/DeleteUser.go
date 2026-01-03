package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteUserHandler удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору
// @Tags users
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {integer} int "ID удалённой задачи"
// @Failure 400 {string} string "Неверный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /users/{id} [delete]

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
