package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuthorize(t *testing.T) {
	service := NewAuthorizationService(&server.Config{
		JWT: struct {
			Secret string
			TTL    time.Duration
		}{
			Secret: "jwtsecret",
			TTL:    1 * time.Hour,
		}})

	tokenString, err := service.Authorize(t.Context(), "testuser")
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method")
		}
		return []byte("jwtsecret"), nil
	})
	require.NoError(t, err)

	subject, err := token.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, "testuser", subject)
}
