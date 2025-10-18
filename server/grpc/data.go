package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type updateIn interface {
	HasId() bool
	GetId() int64
}

func updateData[I updateIn, U any](
	ctx context.Context,
	in I,
	mapper func(I) U,
	updater func(context.Context, int64, U) error,
	logger *log.Logger,
) (*empty.Empty, error) {
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	data := mapper(in)

	if err := updater(ctx, in.GetId(), data); err != nil {
		if errors.Is(err, server.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, server.ErrPermissionDenied.Error())
		}
		logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

type removable interface {
	HasId() bool
	GetId() int64
}

func removeData(
	ctx context.Context,
	in removable,
	remove func(context.Context, int64) error,
	logger *log.Logger,
) (*empty.Empty, error) {
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := remove(ctx, in.GetId()); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
