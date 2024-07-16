package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"fliqt/config"
)

func NewS3PresignClient(cfg *config.Config) (*s3.PresignClient, error) {
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

	client := s3.NewFromConfig(s3cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
	})

	return s3.NewPresignClient(client), nil
}
