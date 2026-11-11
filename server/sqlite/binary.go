package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

var _ server.BinaryService = (*BinaryService)(nil)

type BinaryService struct {
	qs             *sqlc.Queries
	converter      converter.DataConverter
	binariesFolder string
}

func NewBinaryService(queries *sqlc.Queries, db *DB, converter converter.DataConverter) *BinaryService {
	return &BinaryService{
		qs:             queries,
		converter:      converter,
		binariesFolder: db.binariesFolder,
	}
}

func (s *BinaryService) Create(ctx context.Context, data server.BinaryData) (int64, time.Time, error) {
	result, err := s.qs.InsertBinary(ctx, s.converter.ConvertToInsertBinary(ctx, data))
	return result.ID, time.Unix(result.UpdatedAt, 0), unwrapInsertError(err)
}

func (s *BinaryService) Upload(ctx context.Context, id int64, content io.Reader) error {
	asset, err := os.Create(s.getBinaryAssetPath(id))
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	defer asset.Close()

	if _, err = io.CopyBuffer(asset, content, make([]byte, 1024*1024)); err != nil {
		os.Remove(asset.Name())
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s *BinaryService) GetAll(ctx context.Context) ([]server.BinaryData, error) {
	return getAllData(ctx, s.qs.SelectBinaries, s.converter.ConvertToBinaryDataSlice)
}

func (s *BinaryService) Update(ctx context.Context, id int64, data server.BinaryDataUpdate) (time.Time, error) {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, server.ErrDataNotFound
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateBinary(binary)
	s.converter.ConvertToUpdateBinaryUpdate(data, &params)

	updatedAt, err := s.qs.UpdateBinary(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, server.ErrDataNotFound
		}
		return time.Time{}, fmt.Errorf("update: %w", err)
	}
	return time.Unix(updatedAt, 0), nil
}

func (s *BinaryService) Remove(ctx context.Context, id int64) error {
	if err := removeData(ctx, s.qs.DeleteBinary, id); err != nil {
		return err
	}

	if err := os.Remove(s.getBinaryAssetPath(id)); err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	return nil
}

func (s *BinaryService) Open(ctx context.Context, id int64) (io.ReadSeekCloser, int64, error) {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, server.ErrDataNotFound
	}
	if err != nil {
		return nil, 0, fmt.Errorf("open: select: %w", err)
	}

	file, err := os.Open(s.getBinaryAssetPath(id))
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0, server.ErrUploadIncomplete
	}
	if err != nil {
		return nil, 0, fmt.Errorf("open: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("open: stat: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() != binary.Size {
		file.Close()
		return nil, 0, fmt.Errorf("open: invalid completed file or size mismatch")
	}

	return file, info.Size(), nil
}

func (s *BinaryService) WriteChunk(
	ctx context.Context, id int64, offset int64, data []byte, finish bool,
) (committedSize int64, err error) {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return 0, server.ErrDataNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("write chunk: select: %w", err)
	}
	if offset < 0 {
		return 0, server.ErrInvalidOffset
	}
	if binary.Size < 0 {
		return 0, fmt.Errorf("write chunk: negative expected size")
	}

	if _, err := os.Stat(s.getBinaryAssetPath(id)); err == nil {
		return 0, fmt.Errorf("write chunk: upload already complete")
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("write chunk: stat completed file: %w", err)
	}

	file, err := os.OpenFile(s.getIncompleteBinaryAssetPath(id), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return 0, fmt.Errorf("write chunk: open: %w", err)
	}
	closed := false
	defer func() {
		if !closed {
			if closeErr := file.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("write chunk: close: %w", closeErr)
			}
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("write chunk: stat: %w", err)
	}
	committedSize = info.Size()
	if info.Size() != offset {
		return committedSize, server.ErrInvalidOffset
	}
	if offset > binary.Size || int64(len(data)) > binary.Size-offset {
		return committedSize, fmt.Errorf("write chunk: data exceeds expected size %d", binary.Size)
	}
	if err := ctx.Err(); err != nil {
		return committedSize, err
	}

	n, err := file.WriteAt(data, offset)
	committedSize += int64(n)
	if err != nil {
		return committedSize, fmt.Errorf("write chunk: write: %w", err)
	}
	if n != len(data) {
		return committedSize, fmt.Errorf("write chunk: write: %w", io.ErrShortWrite)
	}
	if err := file.Sync(); err != nil {
		return committedSize, fmt.Errorf("write chunk: sync: %w", err)
	}
	if !finish {
		return committedSize, nil
	}
	if committedSize != binary.Size {
		return committedSize, server.ErrUploadIncomplete
	}

	err = file.Close()
	closed = true
	if err != nil {
		return committedSize, fmt.Errorf("write chunk: close: %w", err)
	}
	if err := os.Rename(s.getIncompleteBinaryAssetPath(id), s.getBinaryAssetPath(id)); err != nil {
		return committedSize, fmt.Errorf("write chunk: rename: %w", err)
	}
	return committedSize, nil
}

func (s *BinaryService) WriteStatus(ctx context.Context, id int64) (committedSize int64, complete bool, err error) {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, server.ErrDataNotFound
	}
	if err != nil {
		return 0, false, fmt.Errorf("write status: select: %w", err)
	}

	finalPath := s.getBinaryAssetPath(id)
	for _, path := range []string{finalPath, s.getIncompleteBinaryAssetPath(id), finalPath} {
		info, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return 0, false, fmt.Errorf("write status: stat: %w", err)
		}

		complete := path == finalPath
		if !info.Mode().IsRegular() || info.Size() > binary.Size || (complete && info.Size() != binary.Size) {
			return 0, false, fmt.Errorf("write status: invalid file or size mismatch")
		}
		return info.Size(), complete, nil
	}
	return 0, false, server.ErrDataNotFound
}

func (s *BinaryService) getBinaryAssetPath(id int64) string {
	return filepath.Join(s.binariesFolder, fmt.Sprintf("%d", id))
}

func (s *BinaryService) getIncompleteBinaryAssetPath(id int64) string {
	return filepath.Join(s.binariesFolder, fmt.Sprintf("%d.part", id))
}
