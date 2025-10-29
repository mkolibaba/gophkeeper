package sqlite

import (
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoginCreate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateLogin(t, "app1", "login1", "123", "alice")
	mustCreateLogin(t, "app2", "login2", "123", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM login")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewLoginService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.LoginData{
			Name:  "app3",
			Login: "login3",
		})
		require.NoError(t, err)
	})
	t.Run("user_not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		err := srv.Create(ctx, server.LoginData{
			Name:  "app3",
			Login: "login3",
		})
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
	t.Run("existing_name", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.LoginData{
			Name:     "app1",
			Login:    "login1",
			Password: "123",
		})
		require.NoError(t, err)
	})
}

func TestLoginGetAll(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	mustCreateLogin(t, "app1", "login1", "123", "alice")
	mustCreateLogin(t, "app2", "login2", "123", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM login")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewLoginService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		logins, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, logins, 2)
	})
	t.Run("no_rows", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		logins, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, logins, 0)
	})
	t.Run("non_existent_user", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		logins, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, logins, 0)
	})
}

func TestLoginUpdate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	login1ID := mustCreateLogin(t, "app1", "login1", "123", "alice")
	mustCreateLogin(t, "app2", "login2", "123", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM login")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewLoginService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		password := "superpassword"
		err := srv.Update(ctx, login1ID, server.LoginDataUpdate{
			Password: &password,
		})
		require.NoError(t, err)

		updatedLogin, err := queries.SelectLogin(ctx, login1ID, "alice")
		require.NoError(t, err)
		require.Equal(t, password, *updatedLogin.Password)
		require.Equal(t, "app1", updatedLogin.Name)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		password := "superpassword"
		err := srv.Update(ctx, -100, server.LoginDataUpdate{
			Password: &password,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		password := "superpassword"
		err := srv.Update(ctx, login1ID, server.LoginDataUpdate{
			Password: &password,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func TestLoginDelete(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	login1ID := mustCreateLogin(t, "app1", "login1", "123", "alice")
	mustCreateLogin(t, "app2", "login2", "123", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM login")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewLoginService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		note3ID := mustCreateLogin(t, "app3", "login3", "", "alice")

		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Remove(ctx, note3ID)
		require.NoError(t, err)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Remove(ctx, -100)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		err := srv.Remove(ctx, login1ID)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func mustCreateLogin(t *testing.T, name string, login string, password string, user string) int64 {
	id, err := queries.InsertLogin(t.Context(), sqlc.InsertLoginParams{
		Name:     name,
		Login:    login,
		Password: &password,
		User:     user,
	})
	require.NoError(t, err)
	return id
}
