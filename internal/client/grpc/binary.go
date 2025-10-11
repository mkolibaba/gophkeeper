package grpc

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
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
	file, err := os.Open(data.Filename)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	stream, err := b.client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	buffer := make([]byte, 64*1024)
	var idx int32

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("save: %w", err)
		}

		var chunk pb.FileChunk
		chunk.SetData(buffer[:n])
		chunk.SetIndex(idx)

		var in pb.SaveBinaryRequest
		in.SetChunk(&chunk)
		in.SetName(data.Name)
		in.SetFilename(data.Filename[strings.LastIndex(data.Filename, "/")+1:])
		in.SetSize(fileInfo.Size())
		in.SetNotes(data.Notes)

		if err := stream.Send(&in); err != nil {
			return fmt.Errorf("save: %w", err)
		}

		idx++
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
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
			Filename: b.GetFilename(),
			Size:     b.GetSize(),
			Notes:    b.GetNotes(),
		})
	}
	return binaries, nil
}

func (b *BinaryService) Download(ctx context.Context, name string) error {
	var in pb.DownloadBinaryRequest
	in.SetName(name)

	stream, err := b.client.Download(ctx, &in)
	if err != nil {
		return err
	}

	var file *os.File
	//var totalSize int64
	var receivedBytes int64

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if file == nil {
			file, err = os.Create(chunk.GetFilename())
			if err != nil {
				panic(err)
			}
			defer file.Close()

			//totalSize = chunk.GetTotalSize()
		}

		n, err := file.Write(chunk.GetChunk().GetData())
		if err != nil {
			panic(err)
		}

		receivedBytes += int64(n)

		//if totalSize > 0 {
		//	progress := float64(receivedBytes) / float64(totalSize) * 100
		//}
	}

	return nil
}

func (b *BinaryService) Remove(ctx context.Context, name string) error {
	var in pb.RemoveDataRequest
	in.SetName(name)

	_, err := b.client.Remove(ctx, &in)
	return err
}
