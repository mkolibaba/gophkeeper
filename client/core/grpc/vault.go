package grpc

import (
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"google.golang.org/grpc"
)

func NewVaultServiceClient(conn *grpc.ClientConn) gophkeeperv1.VaultServiceClient {
	return gophkeeperv1.NewVaultServiceClient(conn)
}

var _ client.

type VaultService struct {
	client gophkeeperv1.VaultServiceClient
}

func NewVaultService(client gophkeeperv1.VaultServiceClient) *VaultService {
	return &VaultService{
		
	}
} 