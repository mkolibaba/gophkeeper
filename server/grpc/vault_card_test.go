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

func TestVaultCreateCard(t *testing.T) {
	for _, scenario := range []string{"success", "validation_error", "service_error"} {
		t.Run(scenario, func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			service := &mock.CardServiceMock{
				CreateFunc: func(ctx context.Context, data server.CardData) (int64, time.Time, error) {
					require.Equal(t, t.Context(), ctx)
					require.Equal(t, server.CardData{Name: "entry", Notes: "notes", Number: "4242424242424242", ExpDate: "12/25", CVV: "123", Cardholder: "homer"}, data)
					if scenario == "service_error" {
						return 0, time.Time{}, errors.New("db error")
					}
					return 7, stamp, nil
				},
			}
			f.cards = service
			var in pb.CreateItemRequest
			in.SetName("entry")
			in.SetNotes("notes")
			var payload pb.Card
			payload.SetNumber("4242424242424242")
			payload.SetExpDate("12/25")
			payload.SetCvv("123")
			payload.SetCardholder("homer")
			in.SetCard(&payload)
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
				require.Equal(t, "card/7", out.GetId())
				require.Equal(t, "entry", out.GetName())
				require.Equal(t, "notes", out.GetNotes())
				require.Equal(t, stamp, out.GetUpdatedAt().AsTime())
				require.True(t, out.HasCard())
				require.Equal(t, "4242424242424242", out.GetCard().GetNumber())
				require.Equal(t, "12/25", out.GetCard().GetExpDate())
				require.Equal(t, "123", out.GetCard().GetCvv())
				require.Equal(t, "homer", out.GetCard().GetCardholder())
				require.Len(t, service.CreateCalls(), 1)
			}
		})
	}
}

func TestVaultUpdateCard(t *testing.T) {
	for _, scenario := range []string{"clear_fields", "common_fields_only", "omitted_fields", "not_found", "service_error"} {
		t.Run(scenario, func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			service := &mock.CardServiceMock{
				UpdateFunc: func(ctx context.Context, id int64, data server.CardDataUpdate) (time.Time, error) {
					require.Equal(t, t.Context(), ctx)
					require.EqualValues(t, 7, id)
					switch scenario {
					case "clear_fields":
						require.Equal(t, server.CardDataUpdate{Name: new(""), Notes: new(""), Number: new(""), ExpDate: new(""), CVV: new(""), Cardholder: new("")}, data)
					case "common_fields_only":
						require.Equal(t, server.CardDataUpdate{Name: new(""), Notes: new("")}, data)
					default:
						require.Equal(t, server.CardDataUpdate{}, data)
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
			f.cards = service
			var in pb.UpdateItemRequest
			in.SetId("card/7")
			if scenario == "clear_fields" || scenario == "common_fields_only" {
				in.SetName("")
				in.SetNotes("")
			}
			if scenario == "clear_fields" {
				var payload pb.Card
				payload.SetNumber("")
				payload.SetExpDate("")
				payload.SetCvv("")
				payload.SetCardholder("")
				in.SetCard(&payload)
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
				require.True(t, proto.Equal(in.GetCard(), out.GetCard()))
			}
			require.Len(t, service.UpdateCalls(), 1)
		})
	}
}

func TestVaultGetCard(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "service_error"}[fail], func(t *testing.T) {
			f := newVaultFixture()
			stamp := time.Unix(1700000000, 0).UTC()
			f.cards.GetAllFunc = func(ctx context.Context) ([]server.CardData, error) {
				require.Equal(t, t.Context(), ctx)
				if fail {
					return nil, errors.New("db error")
				}
				return []server.CardData{server.CardData{ID: 7, Name: "entry", Notes: "notes", Number: "4242424242424242", ExpDate: "12/25", CVV: "123", Cardholder: "homer", UpdatedAt: stamp}}, nil
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
			require.Equal(t, "card/7", out.GetId())
			require.Equal(t, "entry", out.GetName())
			require.Equal(t, "notes", out.GetNotes())
			require.Equal(t, stamp, out.GetUpdatedAt().AsTime())
			require.True(t, out.HasCard())
			require.Equal(t, "4242424242424242", out.GetCard().GetNumber())
			require.Equal(t, "12/25", out.GetCard().GetExpDate())
			require.Equal(t, "123", out.GetCard().GetCvv())
			require.Equal(t, "homer", out.GetCard().GetCardholder())
		})
	}
}
