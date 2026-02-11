package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"smart-todo/utils"
	"time"

	"github.com/jackc/pgx/v5"
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

	//Получаем context текущего запроса
	// В нём хранится служебная информация (userID, таймауты, отмена запроса)
	ctx := r.Context()

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
	userIDFromToken, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// ПРОВЕРЯЕМ ТОКЕН В БД
	// SQL ЗАПРОС: ищем токен в таблице refresh_tokens
	// Проверяем:
	// 1. Существует ли такой токен в БД
	// 2. Не истек ли срок действия (expires_at > NOW())
	// 3. Совпадает ли user_id из токена с user_id в БД
	var userIDFromDB int
	var expiresAt time.Time

	err = db.QueryRow(
		ctx,
		`SELECT user_id, expires_at 
		 FROM refresh_tokens 
		 WHERE token = $1 AND expires_at > NOW()`,
		req.RefreshToken, // $1 - ищем этот токен
	).Scan(&userIDFromDB, &expiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Токен не найден в БД или истек
			// Возможные причины:
			// 1. Токен был удален при logout
			// 2. Токен никогда не был сохранен (поддельный)
			// 3. Срок действия истек (expires_at < NOW())
			http.Error(w, "Refresh token not found or expired", http.StatusUnauthorized)
			return
		}
		// Другая ошибка БД
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// ПРОВЕРКА БЕЗОПАСНОСТИ: сверяем user_id из токена и user_id из БД
	// Это защита от подделки токенов:
	if userIDFromToken != userIDFromDB {
		// user_id не совпадают - токен подделан!
		http.Error(w, "Token user_id mismatch", http.StatusUnauthorized)
		return
	}

	// Генерируем новый access token
	accessToken, err := utils.GenerateAccessToken(userIDFromDB)
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
