package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteTaskHandler удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {integer} int "ID удалённой задачи"
// @Failure 400 {string} string "Неверный ID"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	varsDelete := mux.Vars(r)
	idStrDelete := varsDelete["id"]

	id, err := strconv.Atoi(idStrDelete)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed from delete tasks:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //Говорит, что ответ будет в формате json
	json.NewEncoder(w).Encode(id)                      //Из Go в json
}
