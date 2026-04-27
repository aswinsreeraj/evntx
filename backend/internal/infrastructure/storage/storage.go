package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aswinsreeraj/evntx/pkg/logger"
)

type Uploader interface {
	UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error)
}

var defaultUploader Uploader

func Init() error {
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "local" // Default to local for dev
	}

	if driver == "s3" {
		region := os.Getenv("AWS_REGION")
		bucket := os.Getenv("S3_BUCKET")

		if region == "" || bucket == "" {
			return fmt.Errorf("missing S3 config: AWS_REGION and S3_BUCKET are required for S3 driver")
		}

		var cfg aws.Config
		var err error

		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

		if accessKey != "" && secretKey != "" {
			cfg, err = awsConfig.LoadDefaultConfig(context.Background(),
				awsConfig.WithRegion(region),
				awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
			)
		} else {
			cfg, err = awsConfig.LoadDefaultConfig(context.Background(),
				awsConfig.WithRegion(region),
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
		logger.Log.Info().Str("bucket", bucket).Str("region", region).Msg("S3 storage initialized")

	} else {
		// Initialize local storage
		baseURL := os.Getenv("API_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		baseURL = strings.TrimSuffix(baseURL, "/")

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return fmt.Errorf("failed to create local upload directory: %w", err)
		}

		defaultUploader = &LocalUploader{
			baseURL:   baseURL,
			uploadDir: uploadDir,
		}
		logger.Log.Info().Str("dir", uploadDir).Str("baseURL", baseURL).Msg("Local storage initialized")
	}

	return nil
}

func UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error) {
	if defaultUploader == nil {
		return "", fmt.Errorf("storage uploader not initialized")
	}
	return defaultUploader.UploadFile(ctx, key, contentType, body)
}

// S3Uploader implementation
type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
}

func (u *S3Uploader) UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Log.Error().Err(err).Str("bucket", u.bucket).Str("key", key).Msg("S3 upload failed")
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	logger.Log.Info().Str("url", url).Msg("S3 upload successful")
	return url, nil
}

// LocalUploader implementation
type LocalUploader struct {
	baseURL   string
	uploadDir string
}

func (u *LocalUploader) UploadFile(ctx context.Context, key string, contentType string, body io.Reader) (string, error) {
	fullPath := filepath.Join(u.uploadDir, key)
	
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log.Error().Err(err).Str("dir", dir).Msg("Failed to create local directory for upload")
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		logger.Log.Error().Err(err).Str("path", fullPath).Msg("Failed to create local file")
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, body); err != nil {
		logger.Log.Error().Err(err).Str("path", fullPath).Msg("Failed to write to local file")
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	url := fmt.Sprintf("%s/uploads/%s", u.baseURL, key)
	logger.Log.Info().Str("url", url).Msg("Local upload successful")
	return url, nil
}
