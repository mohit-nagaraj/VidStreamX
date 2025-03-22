package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

	qURL := "https://sqs.ap-south-1.amazonaws.com/254797531501/Temp2Process"

	result, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qURL),
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   60,
		WaitTimeSeconds:     10,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

    for _, message := range result.Messages {
        fmt.Printf("Message ID: %s\n", *message.MessageId)
        fmt.Printf("Message Body: %s\n", *message.Body)

        // Delete the message from the queue
        _, err := svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
            QueueUrl:      aws.String(qURL),
            ReceiptHandle: message.ReceiptHandle,
        })

        if err != nil {
            log.Fatalf("Error deleting message: %v", err)
        }

        fmt.Println("Message Deleted")
    }
}
