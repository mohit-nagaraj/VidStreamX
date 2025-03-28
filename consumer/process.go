package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// LaunchECSTask launches an ECS task with the specified configuration
func LaunchECSTask(ecsClient *ecs.Client, bucketName, objectKey string) error {

	input := &ecs.RunTaskInput{
		Cluster:        aws.String("arn:aws:ecs:ap-south-1:254797531501:cluster/TranscoderCluster2"),
		LaunchType:     types.LaunchTypeFargate,
		TaskDefinition: aws.String("arn:aws:ecs:ap-south-1:254797531501:task-definition/video-transcoder"),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{"sg-0a1579071dfd36a04"},
				Subnets: []string{
					"subnet-03da707404a238d17",
					"subnet-0161716fe8b43e42c",
					"subnet-0c1e570cf781d3601",
					"subnet-0dc25ff970f04321e",
				},
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name: aws.String("vidstreamx"),
					Command: []string{
						"-bucket=" + bucketName,
						"-key=" + objectKey,
					},
				},
			},
		},
	}

	result, err := ecsClient.RunTask(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to run ECS task: %w", err)
	}

	for _, task := range result.Tasks {
		fmt.Printf("Launched ECS task: %s\n", *task.TaskArn)
	}

	return nil
}
