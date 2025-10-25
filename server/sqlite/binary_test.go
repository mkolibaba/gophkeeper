package sqlite

import (
	"bytes"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
		err := srv.Create(ctx, server.ReadableBinaryData{
			BinaryData: server.BinaryData{
				Name:     "text_3",
				Filename: "text3.txt",
				Size:     100,
			},
			DataReader: io.NopCloser(strings.NewReader("hello")),
		})
		require.NoError(t, err)
	})
	t.Run("user_not_found", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "charlie")
		err := srv.Create(ctx, server.ReadableBinaryData{
			BinaryData: server.BinaryData{
				Name:     "text_3",
				Filename: "text3.txt",
				Size:     100,
			},
			DataReader: io.NopCloser(strings.NewReader("hello")),
		})
		require.ErrorIs(t, err, server.ErrUserNotFound)
	})
	t.Run("existing_name", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "alice")
		err := srv.Create(ctx, server.ReadableBinaryData{
			BinaryData: server.BinaryData{
				Name:     "text_3",
				Filename: "text3.txt",
				Size:     100,
			},
			DataReader: io.NopCloser(strings.NewReader("hello")),
		})
		require.NoError(t, err)
	})
}

func TestBinaryGen(t *testing.T) {
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
		binaryData, err := srv.Get(ctx, binaryID)
		require.NoError(t, err)
		require.Equal(t, binaryData.Name, "text")
		content, err := io.ReadAll(binaryData.DataReader)
		require.NoError(t, err)
		require.Equal(t, "content 1", string(content))
	})
	// TODO: почему-то err == nil, поправить
	//t.Run("file_removed", func(t *testing.T) {
	//	id := mustCreateBinary(t, "text2", "text_2.txt", strings.NewReader("content 2"), "alice")
	//	// TODO(minor): имплементация скопирована из сервиса. возможно, не очень хорошо так делать
	//	require.NoError(t, os.Remove(srv.getBinaryAssetPath(id)))
	//
	//	ctx := server.NewContextWithUser(t.Context(), "alice")
	//	_, err := srv.Get(ctx, binaryID)
	//	require.Error(t, err)
	//})
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
		err := srv.Update(ctx, binary1ID, server.BinaryDataUpdate{
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
		err := srv.Update(ctx, -100, server.BinaryDataUpdate{
			Name: &name,
		})
		require.ErrorIs(t, err, server.ErrDataNotFound)
	})
	t.Run("user_not_owner", func(t *testing.T) {
		ctx := server.NewContextWithUser(t.Context(), "bob")
		name := "new name"
		err := srv.Update(ctx, binary1ID, server.BinaryDataUpdate{
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

	id, err := queries.InsertBinary(t.Context(), sqlc.InsertBinaryParams{
		Name:     name,
		Filename: filename,
		Size:     size,
		Notes:    nil,
		User:     user,
	})
	require.NoError(t, err)

	// TODO(minor): имплементация скопирована из сервиса. возможно, не очень хорошо так делать
	path := filepath.Join(db.binariesFolder, fmt.Sprintf("%d", id))
	err = os.WriteFile(path, buf.Bytes(), 0666)
	require.NoError(t, err)

	return id
}

func getPokemonPicture(t *testing.T) io.Reader {
	content, err := os.ReadFile(filepath.Join("testdata", "squirtle_pokemon.png"))
	require.NoError(t, err)
	return bytes.NewReader(content)
}
