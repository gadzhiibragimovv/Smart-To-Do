package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) { //Принимает обычный пароль "123456"
	//возращает хеш или ошибку если bcrypt не смог захешировать
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //bcrypt работает только с байтами пример хеша $2a$10$H2m9ZK7P0s2y...Q9x
	//DefaultCost время хеширования пароля, чем выше cost тем медленее
	return string(hash), err
}

// CheckPassword сравнивает пароль и хеш
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) //error - пароль неверный, nil - пароль верный
}
