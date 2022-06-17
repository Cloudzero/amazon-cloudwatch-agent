
# CloudZero Agent
The CloudZero Agent is used to track Pod performance metrics and send to CloudWatch Log for consumption by the CLoudZero platform. The pod performance metrics are correlated across Kubernetes nodes generating [Container Cost](https://docs.cloudzero.com/docs/container-cost-track).

This code is a fork of the Amazon [CloudWatch Agent](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Install-CloudWatch-Agent.html)
 with performance enhancements, disable expensive container insight metric publication and reduction of the number of metrics collection to match the CloudZero use case.

 The container is built as a multi-arch container supporting arm64 and amd64 chip architectures.

## Overview
The Amazon CloudWatch Agent enables you to do the following:

- Collect pod metrics from Amazon EC2 instances across operating systems.

The agents is deployed via a Helm Chart as a DeamonSet.

Amazon CloudWatch Agent uses the open-source project [telegraf](https://github.com/influxdata/telegraf) as its dependency. It operates by starting a telegraf agent with some original plugins and some customized plugins.

### Setup
The installation of CloudZero agent is documented in CloudZero [public documentation](<https://docs.cloudzero.com/docs/installation-of-cloudzero-integration>)

The agent being a fork of the Amazon Cloudwatch agent can be permission several ways. The chart prerequisite use a broad solution that satisfy most firms needs.  More restrictive permission can follow the AWS IAM Roles [Configuring IAM Roles](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/create-iam-roles-for-cloudwatch-agent.html)

### Troubleshooting
The installation instruction has a [validation](<https://docs.cloudzero.com/docs/installation-of-cloudzero-integration#validation>) section.

In the field the following have been the main cause of issues:

- Agent not having permission to write to the CloudWatch LogGroup
- Taints and Toleration preventing the agent from running.
- The account the Kubernetes cluster not connected to the CloudZero platform as "resource account"

### Release Branch
Being a fork and making changes that will probably not accepted by the upstream Amazon Open Source repository, CloudZero is using the branch [cloudzero-optimized](https://github.com/Cloudzero/amazon-cloudwatch-agent/tree/cloudzero-optimized) as the branch to build and release from. Other branches are used to stay current with the upstream fork and to pull into the cloudzero-optimized branch.

To build your own version of this again, review the github action workflow named [multiarch-build.yaml](https://github.com/Cloudzero/amazon-cloudwatch-agent/blob/cloudzero-optimized/.github/workflows/multiarch-build.yaml)
