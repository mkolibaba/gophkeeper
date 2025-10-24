package sqlite

import (
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoteCreate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateNote(t, "note1", "some text", "alice")
	mustCreateNote(t, "note2", "another text", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM note")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewNoteService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.NoteData{
			Name: "note3",
			Text: "third note text",
		})
		require.NoError(t, err)
	})
	t.Run("user_not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		err := srv.Create(ctx, server.NoteData{
			Name: "note3",
			Text: "third note text",
		})
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
	t.Run("existing_name", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.NoteData{
			Name: "note1",
			Text: "some text",
		})
		require.NoError(t, err)
	})
}

func TestNoteGetAll(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	mustCreateNote(t, "note1", "some text", "alice")
	mustCreateNote(t, "note2", "another text", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM note")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewNoteService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		notes, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, notes, 2)
	})
	t.Run("no_rows", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		notes, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, notes, 0)
	})
	t.Run("non_existent_user", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		notes, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, notes, 0)
	})
}

func TestNoteUpdate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	note1ID := mustCreateNote(t, "note1", "some text", "alice")
	mustCreateNote(t, "note2", "another text", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM note")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewNoteService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		text := "brand new text"
		err := srv.Update(ctx, note1ID, server.NoteDataUpdate{
			Text: &text,
		})
		require.NoError(t, err)

		updatedNote, err := queries.SelectNote(ctx, note1ID, "alice")
		require.NoError(t, err)
		require.Equal(t, text, *updatedNote.Text)
		require.Equal(t, "note1", updatedNote.Name)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		text := "brand new text"
		err := srv.Update(ctx, -100, server.NoteDataUpdate{
			Text: &text,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		text := "brand new text"
		err := srv.Update(ctx, note1ID, server.NoteDataUpdate{
			Text: &text,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func TestNoteDelete(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	note1ID := mustCreateNote(t, "note1", "some text", "alice")
	mustCreateNote(t, "note2", "another text", "alice")

	t.Cleanup(func() {
		db.db.Exec("DELETE FROM note")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewNoteService(queries, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		note3ID := mustCreateNote(t, "note3", "some text", "alice")

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
		err := srv.Remove(ctx, note1ID)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func mustCreateNote(t *testing.T, name string, text string, user string) int64 {
	id, err := queries.InsertNote(t.Context(), sqlc.InsertNoteParams{
		Name: name,
		Text: &text,
		User: user,
	})
	require.NoError(t, err)
	return id
}
