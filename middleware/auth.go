package middleware

import (
	"net/http"
	"smart-todo/utils"
	"strings"
)

// AuthMiddleware — middleware авторизации
// Либо:
//   - пропускает запрос дальше
//   - либо останавливает его с ошибкой 401
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")

		// Если заголовка нет — запрос не авторизован
		if authHeader == "" {
			// Подсказываем клиенту, что используется Bearer-аутентификация
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`)
			// Отдаём 401 Unauthorized
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`)
			http.Error(w, "Invalid authorization format. Expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		// Убираем "Bearer " и получаем сам токен
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверяем токен и получаем userID
		// TokenVerification:
		//   - валидирует JWT
		//   - возвращает userID
		userID, err := utils.TokenVerification(token)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Кладём userID в context запроса
		// Теперь userID будет доступен во всех handler'ах ниже
		ctx := SetUserID(r.Context(), userID)

		// Создаём новый request с обновлённым context
		r = r.WithContext(ctx)

		// Если всё успешно — передаём управление следующему handler'у
		next.ServeHTTP(w, r)
	})
}
