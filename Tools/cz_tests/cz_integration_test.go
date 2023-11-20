package cz_tests

import (
	"flag"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	logGroupNameFlag = flag.String("logGroupName", "", "Name of the CloudWatch log group")
	regionFlag       = flag.String("region", "", "AWS region")
)

func getCloudWatchLogEvents(logGroupName, region string) ([]*cloudwatchlogs.OutputLogEvent, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	client := cloudwatchlogs.New(sess)

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		Limit:         aws.Int64(1),
		LogStreamName: nil,
	}

	resp, err := client.GetLogEvents(input)
	if err != nil {
		return nil, err
	}

	return resp.Events, nil
}

func TestValidateCloudWatchLogEvent(t *testing.T) {
	flag.Parse()

	if *logGroupNameFlag == "" {
		t.Fatal("Please provide the CloudWatch log group name using the -logGroupName flag")
	}

	region := *regionFlag
	if region == "" {
		// Set your default AWS region here if not using the -region flag
		region = os.Getenv("AWS_REGION")
		if region == "" {
			t.Fatal("Please provide the AWS region using the -region flag or AWS_REGION environment variable")
		}
	}

	logEvents, err := getCloudWatchLogEvents(*logGroupNameFlag, region)

	if err != nil {
		t.Fatalf("Error getting CloudWatch log events: %v", err)
	}

	if len(logEvents) == 0 {
		t.Fatal("No log events found.")
	}

	// Validate the content of the log event, adjust this part based on your log format
	expectedContent := "Expected Log Content"
	actualContent := *logEvents[0].Message

	if actualContent != expectedContent {
		t.Fatalf("Log content does not match. Expected: %s, Actual: %s", expectedContent, actualContent)
	}

	t.Logf("Log Event Content: %s", actualContent)
}
