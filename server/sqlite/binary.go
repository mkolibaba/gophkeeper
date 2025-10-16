package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"io"
	"os"
	"path/filepath"
)

type BinaryService struct {
	qs             *sqlc.Queries
	binariesFolder string
}

func NewBinaryService(queries *sqlc.Queries, db *DB) *BinaryService {
	return &BinaryService{
		qs:             queries,
		binariesFolder: db.binariesFolder,
	}
}

func (s *BinaryService) Create(ctx context.Context, data server.ReadableBinaryData) error {
	id, err := s.qs.InsertBinary(ctx, sqlc.InsertBinaryParams{
		Name:     data.Name,
		Filename: data.Filename,
		Size:     data.Size,
		Notes:    stringOrNull(data.Notes),
		User:     server.UserFromContext(ctx),
	})

	if err != nil {
		return unwrapInsertError(err)
	}

	asset, err := os.Create(s.getBinaryAssetPath(id))
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	defer asset.Close()

	if _, err = io.CopyBuffer(asset, data.DataReader, make([]byte, 1024*1024)); err != nil {
		os.Remove(asset.Name())
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (s *BinaryService) Get(ctx context.Context, id int64) (*server.ReadableBinaryData, error) {
	binary, err := s.qs.SelectBinary(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrDataNotFound
		}
		return nil, fmt.Errorf("get: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, binary); err != nil {
		return nil, err
	}

	file, err := os.Open(s.getBinaryAssetPath(id))
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &server.ReadableBinaryData{
		BinaryData: server.BinaryData{
			ID:       binary.ID,
			Name:     binary.Name,
			Filename: binary.Filename,
			Notes:    stringOrEmpty(binary.Notes),
			Size:     binary.Size,
			User:     binary.User,
		},
		DataReader: file,
	}, nil
}

func (s *BinaryService) GetAll(ctx context.Context) ([]server.BinaryData, error) {
	user := server.UserFromContext(ctx)

	binaries, err := s.qs.SelectBinaries(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.BinaryData
	for _, binary := range binaries {
		result = append(result, server.BinaryData{
			ID:       binary.ID,
			Name:     binary.Name,
			Filename: binary.Filename,
			Size:     binary.Size,
			Notes:    stringOrEmpty(binary.Notes),
			User:     user,
		})
	}

	return result, nil
}

func (s *BinaryService) Update(ctx context.Context, id int64, data server.BinaryDataUpdate) error {
	binary, err := s.qs.SelectBinary(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, binary); err != nil {
		return err
	}

	params := sqlc.UpdateBinaryParams{
		Name:  binary.Name,
		Notes: binary.Notes,
		ID:    id,
	}

	if data.Name != nil {
		params.Name = *data.Name
	}
	if data.Notes == nil {
		params.Notes = data.Notes
	}

	n, err := s.qs.UpdateBinary(ctx, params)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("update: no rows")
	}
	return nil
}

func (s *BinaryService) Remove(ctx context.Context, id int64) error {
	if err := removeData(ctx, s.qs.SelectBinaryUser, s.qs.DeleteBinary, id); err != nil {
		return err
	}

	if err := os.Remove(s.getBinaryAssetPath(id)); err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	return nil
}

func (s *BinaryService) getBinaryAssetPath(id int64) string {
	return filepath.Join(s.binariesFolder, fmt.Sprintf("%d", id))
}
