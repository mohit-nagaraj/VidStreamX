package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" || AWS_REGION == "" {
		log.Fatalf("AWS_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set in the environment.")
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
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// SQS service client
	svc := sqs.NewFromConfig(cfg)

	// ECS client
	ecsClient := ecs.NewFromConfig(cfg)

	qURL := "https://sqs.ap-south-1.amazonaws.com/254797531501/Temp2Process"

	// Continuous polling loop
	for {
		result, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(qURL),
			MaxNumberOfMessages: 1,
			VisibilityTimeout:   60,
			WaitTimeSeconds:     10,
		})

		if err != nil {
			fmt.Println("Error receiving messages:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(result.Messages) == 0 {
			fmt.Println("No messages received")
			continue
		}

		for _, message := range result.Messages {
			var s3Event events.S3Event
			fmt.Printf("Message Body: %s\n", *message.Body)
			errBdy := json.Unmarshal([]byte(*message.Body), &s3Event)
			if errBdy != nil {
				fmt.Println("Error parsing message body:", err)
				continue
			}
			if len(s3Event.Records) > 0 && s3Event.Records[0].EventName == "s3:TestEvent" {
				fmt.Println("Skipping test event")
				continue
			}

			fmt.Printf("Message ID: %s\n", *message.MessageId)

			for _, record := range s3Event.Records {
				bucketName := record.S3.Bucket.Name
				objectKey := record.S3.Object.Key
				fmt.Printf("Processing file %s from bucket %s\n", objectKey, bucketName)

				// launch ecs task: use ecs create task on ui to all values there n put it in here
				if err := LaunchECSTask(ecsClient, bucketName, objectKey); err != nil {
					fmt.Printf("Failed to launch ECS task: %v\n", err)
				}
			}

			// Delete the message from the queue after processing
			_, err = svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(qURL),
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				fmt.Println("Error deleting message:", err)
				continue
			}
			fmt.Println("Message deleted successfully")
		}
	}
}
