package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/mohit-nagaraj/vidstreamx/transcoder/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bucketName := flag.String("bucket", "", "S3 bucket name")
	objectKey := flag.String("key", "", "S3 object key")
	flag.Parse()

	if *bucketName == "" || *objectKey == "" {
		log.Fatal("Bucket name and object key must be provided")
	}

	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" || AWS_REGION == "" {
		log.Fatal("AWS_REGION, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY must be set in the environment")
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		)),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	videoPath := "videos/original-video.mp4"
	if err := utils.DownloadVideoFromS3(s3Client, *bucketName, *objectKey, videoPath); err != nil {
		log.Fatalf("Failed to download video: %v", err)
	}

	resolutions := []struct {
		width  int
		height int
	}{
		{1920, 1080},
		{1280, 720},
		{854, 480},
		{640, 360},
		{426, 240},
	}

	for _, res := range resolutions {
		outputPath := fmt.Sprintf("formatted/%dp.mp4", res.height)
		if err := utils.TranscodeVideo(videoPath, outputPath, res.width, res.height); err != nil {
			log.Printf("Failed to transcode video to %dp: %v", res.height, err)
			continue
		}
	}
	fmt.Println("All videos done processing")
}
