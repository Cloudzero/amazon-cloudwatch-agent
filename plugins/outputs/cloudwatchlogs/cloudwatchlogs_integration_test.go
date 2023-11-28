package cloudwatchlogs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
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

func TestLogGroupExistsIntegration(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}
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
