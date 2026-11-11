package sqlite

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"github.com/stretchr/testify/require"
)

func TestBinaryCreate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateBinary(t, "text", "text_1.txt", strings.NewReader("content 1"), "alice")
	mustCreateBinary(t, "pic", "squirtle_pokemon.png", getPokemonPicture(t), "alice")

	t.Cleanup(func() {
		// TODO: подчищать файлы
		db.db.Exec("DELETE FROM binary")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewBinaryService(queries, db, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		_, _, err := srv.Create(ctx, server.BinaryData{
			Name:     "text_3",
			Filename: "text3.txt",
			Size:     100,
		})
		require.NoError(t, err)
	})
	t.Run("user_not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		_, _, err := srv.Create(ctx, server.BinaryData{
			Name:     "text_3",
			Filename: "text3.txt",
			Size:     100,
		})
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
	t.Run("existing_name", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		_, _, err := srv.Create(ctx, server.BinaryData{
			Name:     "text_3",
			Filename: "text3.txt",
			Size:     100,
		})
		require.NoError(t, err)
	})
}

func TestBinaryGet(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	binaryID := mustCreateBinary(t, "text", "text_1.txt", strings.NewReader("content 1"), "alice")

	t.Cleanup(func() {
		// TODO: подчищать файлы
		db.db.Exec("DELETE FROM binary")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewBinaryService(queries, db, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		file, size, err := srv.Open(ctx, binaryID)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, file.Close()) })
		require.EqualValues(t, len("content 1"), size)
		content, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Equal(t, "content 1", string(content))
		_, err = file.Seek(0, io.SeekStart)
		require.NoError(t, err)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		file, size, err := srv.Open(ctx, binaryID)
		require.ErrorIs(t, err, server.ErrDataNotFound)
		require.Nil(t, file)
		require.Zero(t, size)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		file, size, err := srv.Open(ctx, -100)
		require.ErrorIs(t, err, server.ErrDataNotFound)
		require.Nil(t, file)
		require.Zero(t, size)
	})
	t.Run("file_removed", func(t *testing.T) {
		id := mustCreateBinary(t, "removed", "removed.txt", strings.NewReader("content 2"), "alice")
		require.NoError(t, os.Remove(srv.getBinaryAssetPath(id)))

		ctx := server.NewContextWithUser(t.Context(), "alice")
		file, size, err := srv.Open(ctx, id)
		require.ErrorIs(t, err, server.ErrUploadIncomplete)
		require.Nil(t, file)
		require.Zero(t, size)
	})
	t.Run("upload_not_started", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		id, _, err := srv.Create(ctx, server.BinaryData{Name: "pending", Filename: "pending.txt", Size: 10})
		require.NoError(t, err)

		file, size, err := srv.Open(ctx, id)
		require.ErrorIs(t, err, server.ErrUploadIncomplete)
		require.Nil(t, file)
		require.Zero(t, size)
	})
	t.Run("upload_incomplete", func(t *testing.T) {
		id := mustCreateBinary(t, "partial", "partial.txt", strings.NewReader("content 3"), "alice")
		require.NoError(t, os.Rename(srv.getBinaryAssetPath(id), srv.getIncompleteBinaryAssetPath(id)))

		ctx := server.NewContextWithUser(t.Context(), "alice")
		file, size, err := srv.Open(ctx, id)
		require.ErrorIs(t, err, server.ErrUploadIncomplete)
		require.Nil(t, file)
		require.Zero(t, size)
	})
	t.Run("size_mismatch", func(t *testing.T) {
		id := mustCreateBinary(t, "damaged", "damaged.txt", strings.NewReader("content 4"), "alice")
		require.NoError(t, os.Truncate(srv.getBinaryAssetPath(id), 1))

		ctx := server.NewContextWithUser(t.Context(), "alice")
		file, size, err := srv.Open(ctx, id)
		require.Error(t, err)
		require.Nil(t, file)
		require.Zero(t, size)
	})
}

func TestBinaryGetAll(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	mustCreateBinary(t, "text", "text_1.txt", strings.NewReader("content 1"), "alice")
	mustCreateBinary(t, "pic", "squirtle_pokemon.png", getPokemonPicture(t), "alice")

	t.Cleanup(func() {
		// TODO: подчищать файлы
		db.db.Exec("DELETE FROM binary")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewBinaryService(queries, db, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		binaries, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, binaries, 2)
	})
	t.Run("no_rows", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		binaries, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, binaries, 0)
	})
	t.Run("non_existent_user", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		binaries, err := srv.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, binaries, 0)
	})
}

func TestBinaryUpdate(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	binary1ID := mustCreateBinary(t, "text", "text_1.txt", strings.NewReader("content 1"), "alice")
	mustCreateBinary(t, "pic", "squirtle_pokemon.png", getPokemonPicture(t), "alice")

	t.Cleanup(func() {
		// TODO: подчищать файлы
		db.db.Exec("DELETE FROM binary")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewBinaryService(queries, db, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		name := "new name"
		_, err := srv.Update(ctx, binary1ID, server.BinaryDataUpdate{
			Name: &name,
		})
		require.NoError(t, err)

		updatedBinary, err := queries.SelectBinary(ctx, binary1ID, "alice")
		require.NoError(t, err)
		require.Equal(t, name, updatedBinary.Name)
		require.Equal(t, "text_1.txt", updatedBinary.Filename)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		name := "new name"
		_, err := srv.Update(ctx, -100, server.BinaryDataUpdate{
			Name: &name,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		name := "new name"
		_, err := srv.Update(ctx, binary1ID, server.BinaryDataUpdate{
			Name: &name,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func TestBinaryDelete(t *testing.T) {
	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
	binary1ID := mustCreateBinary(t, "text", "text_1.txt", strings.NewReader("content 1"), "alice")
	mustCreateBinary(t, "pic", "squirtle_pokemon.png", getPokemonPicture(t), "alice")

	t.Cleanup(func() {
		// TODO: подчищать файлы
		db.db.Exec("DELETE FROM binary")
		db.db.Exec("DELETE FROM user")
	})

	srv := NewBinaryService(queries, db, NewDataConverter())

	t.Run("success", func(t *testing.T) {
		binary3ID := mustCreateBinary(t, "text_2", "text_2.txt", strings.NewReader("content 2"), "alice")

		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Remove(ctx, binary3ID)
		require.NoError(t, err)
	})
	t.Run("not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Remove(ctx, -100)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		err := srv.Remove(ctx, binary1ID)
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
}

func mustCreateBinary(t *testing.T, name string, filename string, content io.Reader, user string) int64 {
	var buf bytes.Buffer
	size, err := io.Copy(&buf, content)
	require.NoError(t, err)

	result, err := queries.InsertBinary(t.Context(), sqlc.InsertBinaryParams{
		Name:     name,
		Filename: filename,
		Size:     size,
		Notes:    nil,
		User:     user,
	})
	require.NoError(t, err)

	// TODO(minor): имплементация скопирована из сервиса. возможно, не очень хорошо так делать
	path := filepath.Join(db.binariesFolder, fmt.Sprintf("%d", result.ID))
	err = os.WriteFile(path, buf.Bytes(), 0666)
	require.NoError(t, err)

	return result.ID
}

func getPokemonPicture(t *testing.T) io.Reader {
	content, err := os.ReadFile(filepath.Join("testdata", "squirtle_pokemon.png"))
	require.NoError(t, err)
	return bytes.NewReader(content)
}
