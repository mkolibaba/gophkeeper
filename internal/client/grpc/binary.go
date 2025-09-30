package grpc

import (
	"context"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
)

type BinaryService struct {
	client pb.DataServiceClient
}

func NewBinaryService(client pb.DataServiceClient) *BinaryService {
	return &BinaryService{
		client: client,
	}
}

func (b *BinaryService) Save(ctx context.Context, user string, data client.BinaryData) error {
	var binary pb.Binary
	binary.SetName(data.Name)
	binary.SetData(data.Bytes)
	binary.SetMetadata(data.Metadata)

	var in pb.SaveDataRequest
	in.SetUser(user)
	in.SetBinary(&binary)

	_, err := b.client.Save(ctx, &in)
	return err
}

func (b *BinaryService) GetAll(ctx context.Context, user string) ([]client.BinaryData, error) {
	var in pb.GetAllDataRequest
	in.SetDataType(pb.DataType_BINARY)
	in.SetUser(user)

	result, err := b.client.GetAll(ctx, &in)
	if err != nil {
		return nil, err
	}

	var binaries []client.BinaryData
	for _, data := range result.GetData() {
		binary := data.GetBinary()
		binaries = append(binaries, client.BinaryData{
			Name:     binary.GetName(),
			Bytes:    binary.GetData(),
			Metadata: binary.GetMetadata(),
		})
	}
	return binaries, nil
}

func (b *BinaryService) Remove(ctx context.Context, name string, user string) error {
	var in pb.RemoveDataRequest
	in.SetUser(user)
	in.SetName(name)
	in.SetDataType(pb.DataType_BINARY)

	_, err := b.client.Remove(ctx, &in)
	return err
}
