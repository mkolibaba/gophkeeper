package grpc

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
)

type BinaryService struct {
	client pb.BinaryServiceClient
}

func NewBinaryService(conn *grpc.ClientConn) *BinaryService {
	return &BinaryService{
		client: pb.NewBinaryServiceClient(conn),
	}
}

func (b *BinaryService) Save(ctx context.Context, data client.BinaryData) error {
	//var in pb.Binary
	//in.SetName(data.Name)
	//in.SetFileName(data.FileName)
	//in.SetMetadata(data.Metadata)
	//
	//// TODO
	//
	//_, err := b.client.Upload(ctx, &in)
	//return err

	return fmt.Errorf("unimplemented")
}

func (b *BinaryService) GetAll(ctx context.Context) ([]client.BinaryData, error) {
	result, err := b.client.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	var binaries []client.BinaryData
	for _, b := range result.GetResult() {
		binaries = append(binaries, client.BinaryData{
			Name:     b.GetName(),
			FileName: b.GetFileName(),
			Metadata: b.GetMetadata(),
		})
	}
	return binaries, nil
}

func (b *BinaryService) Get(ctx context.Context, name string) (client.BinaryData, error) {
	//var in pb.GetDataRequest
	//in.SetName(name)
	//in.SetUser(user)
	//in.SetDataType(pb.DataType_BINARY)
	//
	//out, err := b.client.Get(ctx, &in)
	//if err != nil {
	//	return client.BinaryData{}, err
	//}
	//
	//binary := out.GetData().GetBinary()
	//return client.BinaryData{
	//	Name:     binary.GetName(),
	//	Bytes:    binary.GetData(),
	//	Metadata: binary.GetMetadata(),
	//}, nil
	return client.BinaryData{}, fmt.Errorf("unimplemented")
}

func (b *BinaryService) Remove(ctx context.Context, name string) error {
	var in pb.RemoveDataRequest
	in.SetName(name)

	_, err := b.client.Remove(ctx, &in)
	return err
}
