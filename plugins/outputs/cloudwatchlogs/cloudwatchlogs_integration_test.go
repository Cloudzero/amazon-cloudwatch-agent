package cloudwatchlogs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"testing"
)

// ClusterConfig represents the structure of the YAML
type ClusterConfig struct {
	APIVersion        string              `yaml:"apiVersion"`
	Kind              string              `yaml:"kind"`
	Metadata          Metadata            `yaml:"metadata"`
	ManagedNodeGroups []ManagedNodeGroups `yaml:"managedNodeGroups"`
}

// Metadata represents the metadata section of the YAML
type Metadata struct {
	Name   string `yaml:"name"`
	Region string `yaml:"region"`
}

// ManagedNodeGroups represents the managedNodeGroups section of the YAML
type ManagedNodeGroups struct {
	Name         string `yaml:"name"`
	InstanceType string `yaml:"instanceType"`
	MinSize      int    `yaml:"minSize"`
	MaxSize      int    `yaml:"maxSize"`
	IAM          IAM    `yaml:"iam"`
}

// IAM represents the iam section of the YAML
type IAM struct {
	AttachPolicyARNs []string `yaml:"attachPolicyARNs"`
}

func getSessionInfo() (*sts.GetCallerIdentityOutput, error) {
	// Load the AWS SDK configuration from the shared config file and environment variables
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create a new STS client to make API calls.
	client := sts.NewFromConfig(cfg)

	// Call the GetCallerIdentity operation to retrieve information about the current session.
	result, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TestSessionInfo(t *testing.T) {
	// Get information about the current AWS session.
	info, err := getSessionInfo()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	assert.Equal(t, "hocking", *info.Account, "Expected String to be Equal")
	// Output the information.
	fmt.Println("AWS Account ID:", *info.Account)
	fmt.Println("User ARN:", *info.Arn)
	fmt.Println("User ID:", *info.UserId)
}

func TestLogGroupExists(t *testing.T) {
	// Get Log Groups
	logGroups, err := GetCloudWatchLogGroups()
	if err != nil {
		log.Fatalf("No Log Groups Returned %v", err)
	}

	//Get the cluster name from the NodeGroups.yaml file
	yamlFile, err := os.ReadFile("../../../NodeGroups.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal YAML data into a struct
	var clusterConfig ClusterConfig
	err = yaml.Unmarshal(yamlFile, &clusterConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	// Test case: Log group exists with cluster name
	exists := logGroupExists("/aws/containerinsights/"+clusterConfig.Metadata.Name+"/performance", logGroups)
	assert.True(t, exists, "Expected log group to exist")

}

func logGroupExists(logGroupName string, logGroups []string) bool {
	for _, lg := range logGroups {
		if lg == logGroupName {
			return true
		}
	}

	return false
}

// GetCloudWatchLogGroups returns a list of CloudWatch log groups
func GetCloudWatchLogGroups() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := cloudwatchlogs.NewFromConfig(cfg)

	var logGroups []string
	input := &cloudwatchlogs.DescribeLogGroupsInput{}

	for {
		resp, err := client.DescribeLogGroups(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		for _, lg := range resp.LogGroups {
			logGroups = append(logGroups, *lg.LogGroupName)
		}

		if resp.NextToken == nil {
			break
		}

		input.NextToken = resp.NextToken
	}

	return logGroups, nil
}
