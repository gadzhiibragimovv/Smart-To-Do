package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/utils"
)

// RefreshTokenHandler обновляет access token
// @Summary Обновить токен
// @Description Получить новый access token используя refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {string} string "Неверный JSON"
// @Failure 401 {string} string "Неверный refresh token"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /auth/refresh [post]
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Проверяем refresh token и получаем userID
	userID, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Генерируем новый access token
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// СОЗДАЕМ ОТВЕТ ДЛЯ КЛИЕНТА
	// Заполняем её данными:
	response := RefreshResponse{
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
