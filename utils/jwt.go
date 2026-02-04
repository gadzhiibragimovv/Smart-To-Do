package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret получает секрет из переменной окружения
// Вызывается каждый раз при генерации/проверке токена
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

func GenerateAccessToken(userID int) (string, error) { //Принимает ID пользователя и итоге возвращает готовый JWT
	claims := jwt.MapClaims{ //Утверждение внутри токена, структура токена
		"user_id": userID,                                //ID пользователя
		"exp":     time.Now().Add(24 * time.Hour).Unix(), //Время жизни токена, сейчас +24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //Создаем токен

	return token.SignedString(getJWTSecret()) //На выходе JWT
}

func TokenVerification(authorization string) (int, error) { //Принимаем строку, возвращаем int или ошибку
	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) { //Парсим jwt, берет токен проверяет токен, проверяет exp
		return getJWTSecret(), nil
	})

	if err != nil { //Если токен подделан, просрочен или строка кривая ошибка
		return 0, err
	}

	if !token.Valid { //Проверяем валидность токена (процесс проверки подлинности и срока действия JSON WEB Token(JWT))
		return 0, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims) //Достаем Claims - данные, которые положили в токен (user_id, exp)
	if !ok {                                   //Если токен не того типа - токен плохой
		return 0, jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(float64) //Достаем user_id, (float64) JWT = JSON, а JSON это числа float64
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	return int(userID), nil //Возвращаем int, токен валиден, пользователь определен
}
