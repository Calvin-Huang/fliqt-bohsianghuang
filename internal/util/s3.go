package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"fliqt/config"
)

func NewS3Client(cfg *config.Config) (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.S3Key,
		cfg.S3Secret,
		"",
	)

	s3cfg, err := s3config.LoadDefaultConfig(
		context.TODO(),
		s3config.WithRegion(cfg.S3Region),
		s3config.WithCredentialsProvider(creds),
		s3config.WithClientLogMode(aws.LogRequestWithBody),
	)

	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(s3cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
	}), nil
}
