package tests

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

// initializeAWSSession creates a new AWS session using the specified profile.
func initializeAWSSession(profile string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{},
		Profile: profile,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// TestSearchLogGroup tests the creation and searching of an AWS CloudWatch log group.
func TestSearchLogGroup(t *testing.T) {
	// Set the AWS profile name
	awsProfile := "profile"

	// Create a new AWS session using the specified profile
	sess, err := initializeAWSSession(awsProfile)
	if err != nil {
		t.Fatalf("Failed to create AWS session: %v", err)
	}

	// Create CloudWatchLogs client
	svc := cloudwatchlogs.New(sess)

	// Specify the log group name to search for
	logGroupName := "your-log-group-name"

	// Search for the log group
	describeLogGroupsInput := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupName),
	}

	describeLogGroupsOutput, err := svc.DescribeLogGroups(describeLogGroupsInput)
	if err != nil {
		t.Fatalf("Error describing log groups: %v", err)
	}

	// Validate that the log group was created
	found := false
	for _, group := range describeLogGroupsOutput.LogGroups {
		if *group.LogGroupName == logGroupName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Log group %s not found", logGroupName)
	}
}
