package cz_tests

import (
	"flag"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	logGroupName = flag.String("logGroupName", "", "Name of the CloudWatch log group")
	region       = flag.String("region", "", "Name of the AWS region")
)

func getCloudWatchLogGroupName(logGroupName, region string) (string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return "", err
	}

	client := cloudwatchlogs.New(sess)

	input := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupName),
		Limit:              aws.Int64(1),
	}

	resp, err := client.DescribeLogGroups(input)
	if err != nil {
		return "", err
	}

	if len(resp.LogGroups) == 0 {
		return "", nil
	}

	return *resp.LogGroups[0].LogGroupName, nil
}

func TestGetCloudWatchLogGroupName(t *testing.T) {
	flag.Parse()

	if *logGroupName == "" {
		t.Fatal("Please provide the CloudWatch log group name using the -logGroupName flag")
	}

	if *region == "" {
		t.Fatal("Please provide the CloudWatch log group name using the -logGroupName flag")
	}
	logGroupName, err := getCloudWatchLogGroupName(*logGroupName, *region)

	if err != nil {
		t.Fatalf("Error getting CloudWatch log group name: %v", err)
	}

	t.Logf("CloudWatch Log Group Name: %s", logGroupName)
}
