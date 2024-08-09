package grpc

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/antibomberman/mego-protos/gen/go/storage"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *serverAPI) PutObject(ctx context.Context, req *pb.PutObjectRequest) (*pb.PutObjectResponse, error) {

	_, err := s.storage.StatObject(ctx, s.cfg.MinioBucket, req.GetFileName(), minio.StatObjectOptions{})
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "File %s already exists", req.GetFileName())
	}

	_, err = s.storage.PutObject(ctx, s.cfg.MinioBucket, req.GetFileName(), bytes.NewReader(req.GetData()), int64(len(req.GetData())), minio.PutObjectOptions{ContentType: req.ContentType})
	if err != nil {
		return nil, err
	}

	return &pb.PutObjectResponse{FileName: req.GetFileName()}, nil
}

func (s *serverAPI) GetObject(ctx context.Context, req *pb.GetObjectRequest) (*pb.GetObjectResponse, error) {
	object, err := s.storage.GetObject(ctx, s.cfg.MinioBucket, req.GetFileName(), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	info, err := object.Stat()
	if err != nil {
		return nil, err
	}

	content := make([]byte, info.Size)
	object.Read(content)

	return &pb.GetObjectResponse{FileName: req.GetFileName(), Data: content, ContentType: info.ContentType}, nil
}

func (s *serverAPI) GetObjectUrl(ctx context.Context, req *pb.GetObjectUrlRequest) (*pb.GetObjectUrlResponse, error) {
	presignedURL, err := s.storage.PresignedGetObject(ctx, s.cfg.MinioBucket, req.GetFileName(), time.Duration(24)*time.Hour, nil)
	fmt.Println("Presigned URL: ", presignedURL.String())

	if err != nil {
		return nil, err
	}

	return &pb.GetObjectUrlResponse{Url: presignedURL.String()}, nil
}

func (s *serverAPI) DeleteObject(ctx context.Context, req *pb.DeleteObjectRequest) (*pb.DeleteObjectResponse, error) {
	err := s.storage.RemoveObject(ctx, s.cfg.MinioBucket, req.GetFileName(), minio.RemoveObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteObjectResponse{Message: "Successfully deleted " + req.GetFileName()}, nil
}
