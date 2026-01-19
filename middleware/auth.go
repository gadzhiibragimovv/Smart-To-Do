package middleware

import (
	"net/http"
	"smart-todo/utils"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler { //Либо пускаю запрос дальше, либо останавливаю ошибкой (401)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //Оборачивает handler и получаем w и r
		authHeader := r.Header.Get("Authorization") //Читаем Authorization
		if authHeader == "" {                       //Если заголовок(Authorization) пустой ошибка
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`)      //Отправляет, чтобы браузер клиенты понимали этот ресурс защищен Bearer-аутентификацией
			http.Error(w, "Authorization header required", http.StatusUnauthorized) //Отправляем 401 если нет заголовка
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") { //Проверка начинается ли авторизация с Bearer-токена
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`) //Отправляет, чтобы браузер клиенты понимали этот ресурс защищен Bearer-аутентификацией
			http.Error(w, "Invalid authorization format. Expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ") //Извлекаем сам токен

		_, err := utils.TokenVerification(token)
		if err != nil { //Проверяем токент через папку utils/jwt.go
			w.Header().Set("WWW-Authenticate", `Bearer realm = "Restrincted"`) //Отправляет, чтобы браузер клиенты понимали этот ресурс защищен Bearer-аутентификацией
			http.Error(w, "Invalid token"+err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) //Если все ок, пускаем дальше
	})
}
