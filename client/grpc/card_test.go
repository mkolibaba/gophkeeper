package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/grpc/mock"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestCardSave(t *testing.T) {
	clientMock := &mock.CardServiceClientForMockingMock{}
	srv := NewCardService(clientMock)

	err := srv.Save(t.Context(), client.CardData{
		Name: "new name",
		CVV:  "01/22",
	})
	require.NoError(t, err)

	cc := clientMock.SaveCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.Equal(t, c.In.GetCvv(), "01/22")
	require.True(t, c.In.HasNumber())
}

func TestCardGetAll(t *testing.T) {
	clientMock := &mock.CardServiceClientForMockingMock{
		GetAllFunc: func(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*gophkeeperv1.GetAllCardsResponse, error) {
			var card1 gophkeeperv1.Card
			card1.SetId(1)
			card1.SetName("name 1")
			card1.SetNumber("112233")
			var out gophkeeperv1.GetAllCardsResponse
			out.SetResult([]*gophkeeperv1.Card{&card1})
			return &out, nil
		},
	}
	srv := NewCardService(clientMock)

	all, err := srv.GetAll(t.Context())
	require.NoError(t, err)

	require.Len(t, all, 1)
	c := all[0]
	require.Equal(t, c.Name, "name 1")
	require.Equal(t, c.Number, "112233")
}

func TestCardUpdate(t *testing.T) {
	clientMock := &mock.CardServiceClientForMockingMock{}
	srv := NewCardService(clientMock)

	name := "new name"
	err := srv.Update(t.Context(), client.CardDataUpdate{
		Name: &name,
	})
	require.NoError(t, err)

	cc := clientMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.False(t, c.In.HasCvv())
	require.False(t, c.In.HasNumber())
}
