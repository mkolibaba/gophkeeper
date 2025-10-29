package server

import "context"

// contextKey представляет внутренний ключ для добавления полей в контекст.
type contextKey int

const (
	// userContextKey - ключ для пользователя.
	userContextKey contextKey = iota + 1
)

// NewContextWithUser возвращает новый контекст с пользователем.
func NewContextWithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext возвращает пользователя из контекста.
func UserFromContext(ctx context.Context) string {
	user, _ := ctx.Value(userContextKey).(string)
	return user
}
