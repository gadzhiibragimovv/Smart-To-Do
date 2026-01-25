package middleware

import "context"

// contextKey — свой тип для ключей context
// Делается специально, чтобы избежать конфликтов ключей
type contextKey string

const (
	// userIDKey — ключ, под которым мы храним userID в context
	// Он приватный для пакета middleware
	userIDKey contextKey = "user_id"
)

// SetUserID кладёт userID в context
// Принимает существующий context и возвращает НОВЫЙ context
func SetUserID(ctx context.Context, userID int) context.Context {
	// WithValue создаёт новый context,
	// в котором по ключу userIDKey лежит userID
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID достаёт userID из context
// Возвращает:
//   - userID (int)
//   - bool — найден ли userID в context и правильного ли он типа
func GetUserID(ctx context.Context) (int, bool) {
	// ctx.Value возвращает any(принимает значения любого типа), поэтому делаем приведение типа к int
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
