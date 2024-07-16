package service

import (
	"context"
	"fliqt/config"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
)

const (
	presignedUrlExpiration = 5 * time.Minute
	presignedUrlCacheTTL   = 30 * time.Second
)

type S3ServiceInterface interface {
	PresignUpload(ctx context.Context, bucket, userID string, objectKey string, contentType string, fileSize int64) (string, error)
	GetPresignDownloadURL(ctx context.Context, bucket, objectKey string) (string, error)
}

type S3Service struct {
	cfg             *config.Config
	redisClient     *redis.Client
	s3PresignClient *s3.PresignClient
}

func NewS3Service(
	cfg *config.Config,
	redisClient *redis.Client,
	s3PresignClient *s3.PresignClient,
) *S3Service {
	return &S3Service{
		cfg,
		redisClient,
		s3PresignClient,
	}
}

func (s *S3Service) PresignUpload(ctx context.Context, bucket string, userID string, objectKey string, contentType string, fileSize int64) (string, error) {
	cacheKey := fmt.Sprintf("upload_tmp:%s/%s", bucket, userID)
	previousPresignedURL, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		return previousPresignedURL, nil
	}

	// Get the presigned URL from S3
	s3Req, err := s.s3PresignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
			// Ensure the loading file is the same size and same type of presinged file.
			ContentType:   aws.String(contentType),
			ContentLength: aws.Int64(fileSize),
		},
		func(po *s3.PresignOptions) {
			po.Expires = presignedUrlExpiration
		},
	)
	if err != nil {
		return "", err
	}

	// Cache the presigned URL, ensure the presigned URL can't be generated too frequently.
	if err := s.redisClient.Set(ctx, cacheKey, s3Req.URL, presignedUrlCacheTTL).Err(); err != nil {
		return s3Req.URL, err
	}

	return s3Req.URL, nil
}

func (s *S3Service) GetPresignDownloadURL(ctx context.Context, bucket string, objectKey string) (string, error) {
	// Get the presigned URL from S3
	s3Req, err := s.s3PresignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
		func(po *s3.PresignOptions) {
			po.Expires = presignedUrlExpiration
		},
	)
	if err != nil {
		return "", err
	}

	return s3Req.URL, nil
}
