package grpc

import (
	"context"
	"errors"
	"io"
	"testing"

	"charm.land/log/v2"
	pb "github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

type vaultFixture struct {
	binaries *mock.BinaryServiceMock
	cards    *mock.CardServiceMock
	logins   *mock.LoginServiceMock
	notes    *mock.NoteServiceMock
}

func newVaultFixture() *vaultFixture {
	return &vaultFixture{
		binaries: &mock.BinaryServiceMock{},
		cards:    &mock.CardServiceMock{},
		logins:   &mock.LoginServiceMock{},
		notes:    &mock.NoteServiceMock{},
	}
}

func (f *vaultFixture) server(t *testing.T) *VaultServiceServer {
	t.Helper()
	return NewVaultServiceServer(f.binaries, f.cards, f.logins, f.notes, newTestValidator(t), log.New(io.Discard))
}

func TestVaultGet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out, err := newVaultFixture().server(t).Get(t.Context(), nil)
		require.NoError(t, err)
		require.Empty(t, out.GetItems())
	})
	t.Run("all_types", func(t *testing.T) {
		f := newVaultFixture()
		f.binaries.GetAllFunc = func(context.Context) ([]server.BinaryData, error) {
			return []server.BinaryData{{ID: 1, Name: "file", Filename: "test.txt"}}, nil
		}
		f.cards.GetAllFunc = func(context.Context) ([]server.CardData, error) {
			return []server.CardData{{ID: 1, Name: "card"}}, nil
		}
		f.logins.GetAllFunc = func(context.Context) ([]server.LoginData, error) {
			return []server.LoginData{{ID: 1, Name: "login"}}, nil
		}
		f.notes.GetAllFunc = func(context.Context) ([]server.NoteData, error) {
			return []server.NoteData{{ID: 1, Name: "note"}}, nil
		}
		out, err := f.server(t).Get(t.Context(), nil)
		require.NoError(t, err)
		require.Len(t, out.GetItems(), 4)
		var ids []string
		for _, item := range out.GetItems() {
			ids = append(ids, item.GetId())
		}
		require.ElementsMatch(t, []string{"binary/1", "card/1", "login/1", "note/1"}, ids)
		require.Len(t, f.binaries.GetAllCalls(), 1)
		require.Len(t, f.cards.GetAllCalls(), 1)
		require.Len(t, f.logins.GetAllCalls(), 1)
		require.Len(t, f.notes.GetAllCalls(), 1)
	})
}

func TestVaultDeleteItem(t *testing.T) {
	for _, typ := range []string{"binary", "card", "login", "note"} {
		for _, scenario := range []string{"success", "not_found", "service_error"} {
			t.Run(typ+"/"+scenario, func(t *testing.T) {
				f := newVaultFixture()
				calls := 0
				remove := func(ctx context.Context, id int64) error {
					calls++
					require.Equal(t, t.Context(), ctx)
					require.EqualValues(t, 7, id)
					switch scenario {
					case "not_found":
						return server.ErrDataNotFound
					case "service_error":
						return errors.New("db error")
					}
					return nil
				}
				f.binaries.RemoveFunc = remove
				f.cards.RemoveFunc = remove
				f.logins.RemoveFunc = remove
				f.notes.RemoveFunc = remove
				var in pb.DeleteItemRequest
				in.SetId(typ + "/7")
				out, err := f.server(t).DeleteItem(t.Context(), &in)
				switch scenario {
				case "not_found":
					requireGrpcError(t, err, codes.NotFound)
				case "service_error":
					requireGrpcError(t, err, codes.Internal)
				default:
					require.NoError(t, err)
					require.NotNil(t, out)
				}
				require.Equal(t, 1, calls)
				for name, count := range map[string]int{
					"binary": len(f.binaries.RemoveCalls()),
					"card":   len(f.cards.RemoveCalls()),
					"login":  len(f.logins.RemoveCalls()),
					"note":   len(f.notes.RemoveCalls()),
				} {
					expected := 0
					if name == typ {
						expected = 1
					}
					require.Equal(t, expected, count, name)
				}
			})
		}
	}
}

func TestVaultInvalidItemID(t *testing.T) {
	for _, id := range []string{"", "7", "note/", "note/invalid", "unknown/7", "note/7/8"} {
		t.Run(id, func(t *testing.T) {
			f := newVaultFixture()
			srv := f.server(t)
			var update pb.UpdateItemRequest
			update.SetId(id)
			_, err := srv.UpdateItem(t.Context(), &update)
			requireGrpcError(t, err, codes.NotFound)
			var remove pb.DeleteItemRequest
			remove.SetId(id)
			_, err = srv.DeleteItem(t.Context(), &remove)
			requireGrpcError(t, err, codes.NotFound)
			require.Empty(t, f.binaries.UpdateCalls())
			require.Empty(t, f.cards.UpdateCalls())
			require.Empty(t, f.logins.UpdateCalls())
			require.Empty(t, f.notes.UpdateCalls())
			require.Empty(t, f.binaries.RemoveCalls())
			require.Empty(t, f.cards.RemoveCalls())
			require.Empty(t, f.logins.RemoveCalls())
			require.Empty(t, f.notes.RemoveCalls())
		})
	}
}
