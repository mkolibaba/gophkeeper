package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func TestVaultCreateLogin(t *testing.T) {
	for _, scenario := range []string{"success", "validation_error", "service_error"} {
		t.Run(scenario, func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			service := &mock.LoginServiceMock{
				CreateFunc: func(ctx context.Context, data server.LoginData) (int64, time.Time, error) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, server.LoginData{Name: "entry", Notes: "notes", Login: "user", Password: "pass", Website: "example.com"}, data)
					if scenario == "service_error" {
						return 0, time.Time{}, errors.New("db error")
					}
					return 7, stamp, nil
				},
			}
			f.logins = service
			var in pb.CreateItemRequest
			in.SetName("entry")
			in.SetNotes("notes")
			var payload pb.Login
			payload.SetLogin("user")
			payload.SetPassword("pass")
			payload.SetWebsite("example.com")
			in.SetLogin(&payload)
			if scenario == "validation_error" {
				in.ClearName()
			}
			out, err := f.server(t).CreateItem(t.Context(), &in)
			switch scenario {
			case "validation_error":
				requireGrpcError(t, err, codes.InvalidArgument)
				require.Empty(t, service.CreateCalls())
			case "service_error":
				requireGrpcError(t, err, codes.Internal)
			default:
				require.NoError(t, err)
				require.Equal(t, "login/7", out.GetId())
				require.Equal(t, "entry", out.GetName())
				require.Equal(t, "notes", out.GetNotes())
				require.Equal(t, stamp, out.GetUpdatedAt().AsTime())
				require.True(t, out.HasLogin())
				require.Equal(t, "user", out.GetLogin().GetLogin())
				require.Equal(t, "pass", out.GetLogin().GetPassword())
				require.Equal(t, "example.com", out.GetLogin().GetWebsite())
				require.Len(t, service.CreateCalls(), 1)
			}
		})
	}
}

func TestVaultUpdateLogin(t *testing.T) {
	for _, scenario := range []string{"clear_fields", "common_fields_only", "omitted_fields", "not_found", "service_error"} {
		t.Run(scenario, func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			service := &mock.LoginServiceMock{
				UpdateFunc: func(ctx context.Context, id int64, data server.LoginDataUpdate) (time.Time, error) {
					require.Equal(t, t.Context(), ctx)
					require.EqualValues(t, 7, id)
					switch scenario {
					case "clear_fields":
						require.Equal(t, server.LoginDataUpdate{Name: new(""), Notes: new(""), Login: new(""), Password: new(""), Website: new("")}, data)
					case "common_fields_only":
						require.Equal(t, server.LoginDataUpdate{Name: new(""), Notes: new("")}, data)
					default:
						require.Equal(t, server.LoginDataUpdate{}, data)
					}
					if scenario == "not_found" {
						return time.Time{}, server.ErrDataNotFound
					}
					if scenario == "service_error" {
						return time.Time{}, errors.New("db error")
					}
					return stamp, nil
				},
			}
			f.logins = service
			var in pb.UpdateItemRequest
			in.SetId("login/7")
			if scenario == "clear_fields" || scenario == "common_fields_only" {
				in.SetName("")
				in.SetNotes("")
			}
			if scenario == "clear_fields" {
				var payload pb.Login
				payload.SetLogin("")
				payload.SetPassword("")
				payload.SetWebsite("")
				in.SetLogin(&payload)
			}
			out, err := f.server(t).UpdateItem(t.Context(), &in)
			switch scenario {
			case "not_found":
				requireGrpcError(t, err, codes.NotFound)
			case "service_error":
				requireGrpcError(t, err, codes.Internal)
			default:
				require.NoError(t, err)
				require.Equal(t, in.GetId(), out.GetId())
				require.Equal(t, stamp, out.GetUpdatedAt().AsTime())
				require.Equal(t, in.HasName(), out.HasName())
				require.Equal(t, in.HasNotes(), out.HasNotes())
				require.True(t, proto.Equal(in.GetLogin(), out.GetLogin()))
			}
			require.Len(t, service.UpdateCalls(), 1)
		})
	}
}

func TestVaultGetLogin(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "service_error"}[fail], func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			f.logins.GetAllFunc = func(ctx context.Context) ([]server.LoginData, error) {
				require.Equal(t, t.Context(), ctx)
				if fail {
					return nil, errors.New("db error")
				}
				return []server.LoginData{server.LoginData{ID: 7, Name: "entry", Notes: "notes", Login: "user", Password: "pass", Website: "example.com", UpdatedAt: stamp}}, nil
			}
			resp, err := f.server(t).Get(t.Context(), nil)
			if fail {
				requireGrpcError(t, err, codes.Internal)
				require.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.Len(t, resp.GetItems(), 1)
			out := resp.GetItems()[0]
			require.Equal(t, "login/7", out.GetId())
			require.Equal(t, "entry", out.GetName())
			require.Equal(t, "notes", out.GetNotes())
			require.Equal(t, stamp, out.GetUpdatedAt().AsTime())
			require.True(t, out.HasLogin())
			require.Equal(t, "user", out.GetLogin().GetLogin())
			require.Equal(t, "pass", out.GetLogin().GetPassword())
			require.Equal(t, "example.com", out.GetLogin().GetWebsite())
		})
	}
}
