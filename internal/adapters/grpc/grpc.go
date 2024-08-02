package grpc

import (
	"antibomberman/mego-storage/internal/config"
	storageGrpc "github.com/antibomberman/mego-protos/gen/go/storage"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
)

type serverAPI struct {
	storageGrpc.UnimplementedStorageServiceServer
	cfg     *config.Config
	storage *minio.Client
}

func Register(gRPC *grpc.Server, cfg *config.Config, client *minio.Client) {
	storageGrpc.RegisterStorageServiceServer(gRPC, &serverAPI{
		cfg:     cfg,
		storage: client,
	})
}
