package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadFileToS3 uploads a file to the specified S3 bucket and key
func UploadFileToS3(client *s3.Client, bucketName, key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to S3: %w", filePath, err)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", filePath, bucketName, key)
	return nil
}

// UploadDirectoryToS3 recursively uploads files from a local directory to an S3 bucket
func UploadDirectoryToS3(client *s3.Client, bucketName, baseKey, localDir string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("localDir %s path %s\n", localDir, path)
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(localDir, path)
			fmt.Printf("relPath %s\n", relPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}

			s3Key := filepath.Join(baseKey, relPath)
			if err := UploadFileToS3(client, bucketName, s3Key, path); err != nil {
				return fmt.Errorf("failed to upload %s: %w", path, err)
			}
		}

		return nil
	})
}
