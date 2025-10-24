package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestCardSave(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.Card
		in.SetName("new card")
		in.SetNumber("4242424242424242")
		in.SetExpDate("12/25")
		in.SetCvv("123")
		in.SetCardholder("homer")

		_, err := srv.Save(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.Card

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("db_error", func(t *testing.T) {
		cardServiceMock := &mock.CardServiceMock{
			CreateFunc: func(_ context.Context, _ server.CardData) error {
				return fmt.Errorf("some error")
			},
		}
		srv := createCardServiceServer(t, cardServiceMock)

		var in gophkeeperv1.Card
		in.SetName("new card")
		in.SetNumber("4242424242424242")
		in.SetExpDate("12/25")
		in.SetCvv("123")
		in.SetCardholder("homer")

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.Internal)
	})
}

func TestCardUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.Card
		in.SetId(1)
		in.SetName("new card name")

		_, err := srv.Update(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.Card
		in.SetName("new card name")

		_, err := srv.Update(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
}

func TestCardRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createCardServiceServer(t, &mock.CardServiceMock{})

		var in gophkeeperv1.RemoveDataRequest

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("not_found", func(t *testing.T) {
		service := &mock.CardServiceMock{
			RemoveFunc: func(ctx context.Context, id int64) error {
				return server.ErrDataNotFound
			},
		}
		srv := createCardServiceServer(t, service)

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.NotFound)
	})
}

func TestCardGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mock.CardServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.CardData, error) {
				return []server.CardData{
						{ID: 1, Name: "card1"},
						{ID: 2, Name: "card2"},
					},
					nil
			},
		}
		srv := createCardServiceServer(t, service)
		resp, err := srv.GetAll(t.Context(), nil)
		require.NoError(t, err)
		require.Len(t, resp.GetResult(), 2)
	})
	t.Run("db_error", func(t *testing.T) {
		service := &mock.CardServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.CardData, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		srv := createCardServiceServer(t, service)
		_, err := srv.GetAll(t.Context(), nil)
		requireGrpcError(t, err, codes.Internal)
	})
}

func createCardServiceServer(t *testing.T, cardService server.CardService) *CardServiceServer {
	return NewCardServiceServer(cardService, newTestValidator(t), log.New(io.Discard))
}
