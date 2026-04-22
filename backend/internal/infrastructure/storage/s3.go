package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

var defaultUploader *S3Uploader

func Init() error {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")

	if region == "" || bucket == "" {
		return fmt.Errorf("missing S3 config: AWS_REGION and S3_BUCKET are required")
	}

	var cfg aws.Config
	var err error

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey != "" && secretKey != "" {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
		)
	}

	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	defaultUploader = &S3Uploader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		region: region,
	}
	return nil
}

func UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error) {
	if defaultUploader == nil {
		return "", fmt.Errorf("S3 uploader not initialized")
	}
	return defaultUploader.UploadFile(ctx, key, contentType, body)
}

func (u *S3Uploader) UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	return url, nil
}
