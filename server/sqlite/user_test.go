package sqlite

import (
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserGet(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	t.Cleanup(func() {
		db.db.Exec("DELETE FROM user")
	})

	srv := NewUserService(queries)

	t.Run("success", func(t *testing.T) {
		user, err := srv.Get(t.Context(), "alice")
		require.NoError(t, err)
		require.Equal(t, "alice", user.Login)
		require.Equal(t, "123", user.Password)
	})
	t.Run("not_found", func(t *testing.T) {
		_, err := srv.Get(t.Context(), "charlie")
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
}

func TestUserSave(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	t.Cleanup(func() {
		db.db.Exec("DELETE FROM user")
	})

	srv := NewUserService(queries)

	t.Run("success", func(t *testing.T) {
		err := srv.Save(t.Context(), server.User{
			Login:    "charlie",
			Password: "123",
		})
		require.NoError(t, err)

		t.Run("hash", func(t *testing.T) {
			user, cErr := srv.Get(t.Context(), "charlie")
			require.NoError(t, cErr)
			require.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("123")))
		})
	})
	t.Run("duplicate", func(t *testing.T) {
		err := srv.Save(t.Context(), server.User{
			Login:    "alice",
			Password: "12333333",
		})
		require.ErrorIs(t, err, server.ErrUserAlreadyExists)
	})
}

func mustCreateUser(t *testing.T, login string, password string) {
	err := queries.InsertUser(t.Context(), login, password)
	require.NoError(t, err)
}
