package main

import (
	adapter "antibomberman/mego-storage/internal/adapters/grpc"
	"antibomberman/mego-storage/internal/config"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)
	gRPC := grpc.NewServer()

	minioClient, err := initMinio(cfg)
	if err != nil {
		log.Printf("Failed to initialize MinIO client: %v", err)
	}
	l, err := net.Listen("tcp", ":"+cfg.StorageServiceServerPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	adapter.Register(gRPC, cfg, minioClient)

	if err := gRPC.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initMinio(cfg *config.Config) (*minio.Client, error) {

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioRootUser, cfg.MinioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		log.Fatalf("Failed to check if bucket exists: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}
	return minioClient, nil
}
