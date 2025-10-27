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

func TestLoginSave(t *testing.T) {
	clientMock := &mock.LoginServiceClientForMockingMock{}
	srv := NewLoginService(clientMock)

	err := srv.Save(t.Context(), client.LoginData{
		Name:  "new name",
		Login: "new login",
	})
	require.NoError(t, err)

	cc := clientMock.SaveCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.Equal(t, c.In.GetLogin(), "new login")
	require.True(t, c.In.HasPassword())
}

func TestLoginGetAll(t *testing.T) {
	clientMock := &mock.LoginServiceClientForMockingMock{
		GetAllFunc: func(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*gophkeeperv1.GetAllLoginsResponse, error) {
			var login1 gophkeeperv1.Login
			login1.SetId(1)
			login1.SetName("name 1")
			login1.SetLogin("login 1")
			var out gophkeeperv1.GetAllLoginsResponse
			out.SetResult([]*gophkeeperv1.Login{&login1})
			return &out, nil
		},
	}
	srv := NewLoginService(clientMock)

	all, err := srv.GetAll(t.Context())
	require.NoError(t, err)

	require.Len(t, all, 1)
	c := all[0]
	require.Equal(t, c.Name, "name 1")
	require.Equal(t, c.Login, "login 1")
}

func TestLoginUpdate(t *testing.T) {
	clientMock := &mock.LoginServiceClientForMockingMock{}
	srv := NewLoginService(clientMock)

	name := "new name"
	err := srv.Update(t.Context(), client.LoginDataUpdate{
		Name: &name,
	})
	require.NoError(t, err)

	cc := clientMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.False(t, c.In.HasPassword())
	require.False(t, c.In.HasWebsite())
}
