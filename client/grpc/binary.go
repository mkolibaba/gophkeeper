package grpc

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
)

type BinaryService struct {
	client gophkeeperv1.BinaryServiceClient
}

func NewBinaryService(conn *grpc.ClientConn) *BinaryService {
	return &BinaryService{
		client: gophkeeperv1.NewBinaryServiceClient(conn),
	}
}

func (s *BinaryService) Save(ctx context.Context, data client.BinaryData) error {
	file, err := os.Open(data.Filename)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	stream, err := s.client.Upload(ctx)
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

		var chunk gophkeeperv1.FileChunk
		chunk.SetData(buffer[:n])
		chunk.SetIndex(idx)

		var in gophkeeperv1.SaveBinaryRequest
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

func (s *BinaryService) GetAll(ctx context.Context) ([]client.BinaryData, error) {
	result, err := s.client.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	var binaries []client.BinaryData
	for _, b := range result.GetResult() {
		binaries = append(binaries, client.BinaryData{
			ID:       b.GetId(),
			Name:     b.GetName(),
			Filename: b.GetFilename(),
			Size:     b.GetSize(),
			Notes:    b.GetNotes(),
		})
	}
	return binaries, nil
}

func (s *BinaryService) Update(ctx context.Context, data client.BinaryDataUpdate) error {
	var in gophkeeperv1.UpdateBinaryRequest
	in.SetId(data.ID)
	if data.Name != nil {
		in.SetName(*data.Name)
	}
	if data.Notes != nil {
		in.SetNotes(*data.Notes)
	}

	_, err := s.client.Update(ctx, &in)
	return err
}

func (s *BinaryService) Download(ctx context.Context, id int64) error {
	var in gophkeeperv1.DownloadBinaryRequest
	in.SetId(id)

	stream, err := s.client.Download(ctx, &in)
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

func (s *BinaryService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := s.client.Remove(ctx, &in)
	return err
}
