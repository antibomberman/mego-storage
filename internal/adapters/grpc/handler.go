package grpc

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/antibomberman/mego-protos/gen/go/storage"
	"github.com/minio/minio-go/v7"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

func (s *serverAPI) PutObject(ctx context.Context, req *pb.PutObjectRequest) (*pb.PutObjectResponse, error) {
	log.Printf("Received PutObjectRequest: %v", req)
	fileName := generateRandomFileName(s.storage, ctx, s.cfg.MinioBucket, req.GetFileName())

	_, err := s.storage.PutObject(ctx, s.cfg.MinioBucket, fileName, bytes.NewReader(req.GetData()), int64(len(req.GetData())), minio.PutObjectOptions{ContentType: req.ContentType})
	if err != nil {
		return nil, err
	}

	return &pb.PutObjectResponse{FileName: fileName}, nil
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

	if err != nil {
		return nil, err
	}
	info, err := s.storage.StatObject(ctx, s.cfg.MinioBucket, req.GetFileName(), minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	return &pb.GetObjectUrlResponse{
		FileName:    req.GetFileName(),
		ContentType: info.ContentType,
		Url:         presignedURL.String(),
	}, nil
}

func (s *serverAPI) DeleteObject(ctx context.Context, req *pb.DeleteObjectRequest) (*pb.DeleteObjectResponse, error) {
	err := s.storage.RemoveObject(ctx, s.cfg.MinioBucket, req.GetFileName(), minio.RemoveObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteObjectResponse{Message: "Successfully deleted " + req.GetFileName()}, nil
}

func generateRandomFileName(storage *minio.Client, ctx context.Context, bucketName, objectName string) string {

	_, err := storage.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		//exist
		ext := filepath.Ext(objectName)
		name := strings.TrimSuffix(objectName, ext)

		for strings.Contains(name, ".") {
			ext = filepath.Ext(name) + ext
			name = strings.TrimSuffix(name, ext)
		}

		rand.Seed(time.Now().UnixNano())
		randomString := fmt.Sprintf("%s_%d", name, rand.Intn(10000))

		return generateRandomFileName(storage, ctx, bucketName, randomString+ext)
	}

	return objectName

}
