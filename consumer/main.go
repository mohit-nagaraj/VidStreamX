package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // SQS service client
    svc := sqs.NewFromConfig(cfg)

    qURL := "https://sqs.ap-south-1.amazonaws.com/254797531501/Temp2Process"

    result, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(qURL),
        MaxNumberOfMessages: 1,
        VisibilityTimeout:   20,
        WaitTimeSeconds:     0,
    })

    if err != nil {
        fmt.Println("Error", err)
        return
    }

    if len(result.Messages) == 0 {
        fmt.Println("Received no messages")
        return
    }

    // msg to be deleted after read
    _, err = svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(qURL),
        ReceiptHandle: result.Messages[0].ReceiptHandle,
    })

    if err != nil {
        fmt.Println("Delete Error", err)
        return
    }

    fmt.Println("Message Deleted")
}
