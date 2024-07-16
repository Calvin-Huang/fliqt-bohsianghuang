package service

import (
	"context"
	"fliqt/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
)

const (
	expiration = 5 * time.Minute
)

type S3ServiceInterface interface {
	PresignUpload(ctx context.Context, bucket, objectKey string, contentType string, fileSize int64) (string, error)
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

func (s *S3Service) PresignUpload(ctx context.Context, bucket string, objectKey string, contentType string, fileSize int64) (string, error) {
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
			po.Expires = expiration
		},
	)
	if err != nil {
		return "", err
	}

	return s3Req.URL, nil
}
