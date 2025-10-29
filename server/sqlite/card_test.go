package sqlite

import (
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCardCreate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateCard(t, "card1", "223322", "01/30", "alice")
	mustCreateCard(t, "card2", "887788", "03/33", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM card")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewCardService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.CardData{
			Name:   "card3",
			Number: "222333333",
		})
		require.NoError(t, err)
	})
	t.Run("user_not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		err := srv.Create(ctx, server.CardData{
			Name:   "card3",
			Number: "222333333",
		})
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
	t.Run("existing_name", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.CardData{
			Name:    "card1",
			Number:  "223322",
			ExpDate: "01/30",
		})
		require.NoError(t, err)
	})
}

func TestCardGetAll(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	mustCreateCard(t, "card1", "223322", "01/30", "alice")
	mustCreateCard(t, "card2", "887788", "03/33", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM card")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewCardService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		cards, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, cards, 2)
	})
	t.Run("no_rows", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		cards, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, cards, 0)
	})
	t.Run("non_existent_user", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		cards, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, cards, 0)
	})
}

func TestCardUpdate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	card1ID := mustCreateCard(t, "card1", "223322", "01/30", "alice")
	mustCreateCard(t, "card2", "887788", "03/33", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM card")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewCardService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		number := "41110011"
		err := srv.Update(ctx, card1ID, server.CardDataUpdate{
			Number: &number,
		})
		require.NoError(t, err)

		updatedCard, err := queries.SelectCard(ctx, card1ID, "alice")
		require.NoError(t, err)
		require.Equal(t, number, updatedCard.Number)
		require.Equal(t, "card1", updatedCard.Name)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		number := "41110011"
		err := srv.Update(ctx, -100, server.CardDataUpdate{
			Number: &number,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		number := "41110011"
		err := srv.Update(ctx, card1ID, server.CardDataUpdate{
			Number: &number,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func TestCardDelete(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	card1ID := mustCreateCard(t, "card1", "223322", "01/30", "alice")
	mustCreateCard(t, "card2", "887788", "03/33", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM card")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewCardService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		note3ID := mustCreateCard(t, "card3", "4444", "01/31", "alice")

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
		err := srv.Remove(ctx, card1ID)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func mustCreateCard(t *testing.T, name string, number string, expdate string, user string) int64 {
	id, err := queries.InsertCard(t.Context(), sqlc.InsertCardParams{
		Name:       name,
		Number:     number,
		ExpDate:    expdate,
		Cvv:        "123",
		Cardholder: strings.ToUpper(user),
		Notes:      nil,
		User:       user,
	})
	require.NoError(t, err)
	return id
}
