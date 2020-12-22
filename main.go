package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"time"
)

/*
{
	"version": "0",
	"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	"detail-type": "ECS Task State Change",
	"source": "aws.ecs",
	"account": "xxxxxxxxxxxx",
	"time": "2020-12-14T15:35:00Z",
	"region": "ap-northeast-1",
	"resources": [
		"arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:task/hoge-hoge-services/xxxxxxxxxxxxxxxxxxxxxxxx"
	],
	"detail": {
		"attachments": [{
			"id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"type": "eni",
			"status": "ATTACHED",
			"details": [{
					"name": "subnetId",
					"value": "subnet-xxxxxxxxxxxxxxxxx"
				},
				{
					"name": "networkInterfaceId",
					"value": "eni-xxxxxxxxxxxxxxxxx"
				},
				{
					"name": "macAddress",
					"value": "xx:xx:xx:xx:xx:xx"
				},
				{
					"name": "privateDnsName",
					"value": "ip-xx-xxx-xx-xxx.ap-northeast-1.compute.internal"
				},
				{
					"name": "privateIPv4Address",
					"value": "xx.xxx.xx.xxx"
				}
			]
		}],
		"availabilityZone": "ap-northeast-1a",
		"capacityProviderName": "FARGATE_SPOT",
		"clusterArn": "arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:cluster/hoge-hoge-services",
		"containers": [{
			"containerArn": "arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:container/xxxxxxxxxxxxxxxxxxxxxxxx",
			"lastStatus": "PENDING",
			"name": "hoge-hoge",
			"image": "xxxxxxxxxxxx.dkr.ecr.ap-northeast-1.amazonaws.com/hoge/hoge@sha256:xxxxxxxxxxxxxxxxxxxxxxxx",
			"taskArn": "arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:task/hoge-hoge-services/xxxxxxxxxxxxxxxxxxxxxxxx",
			"networkInterfaces": [{
				"attachmentId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				"privateIpv4Address": "xx.xx.xx.xxx"
			}],
			"cpu": "0"
		}],
		"cpu": "512",
		"createdAt": "2020-12-14T15:34:41.824Z",
		"desiredStatus": "RUNNING",
		"group": "service:hoge-hoge-hoge-hoge",
		"launchType": "FARGATE",
		"lastStatus": "PENDING",
		"memory": "1024",
		"overrides": {
			"containerOverrides": [{
				"name": "hoge-hoge"
			}]
		},
		"platformVersion": "1.3.0",
		"startedBy": "ecs-svc/7555234739371397917",
		"taskArn": "arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:task/hoge-hoge-services/xxxxxxxxxxxxxxxxxxxxxxxx",
		"taskDefinitionArn": "arn:aws:ecs:ap-northeast-1:xxxxxxxxxxxx:task-definition/hoge-hoge-hoge-hoge:12",
		"updatedAt": "2020-12-14T15:35:00.23Z",
		"version": 2
	}
}
*/
type MyEvent struct {
	Version    string    `json:"version"`
	ID         string    `json:"id"`
	DetailType string    `json:"detail-type"`
	Source     string    `json:"source"`
	Account    string    `json:"account"`
	Time       time.Time `json:"time"`
	Region     string    `json:"region"`
	Resources  []string  `json:"resources"`
	Detail     struct {
		Attachments []struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Status  string `json:"status"`
			Details []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"details"`
		} `json:"attachments"`
		AvailabilityZone     string `json:"availabilityZone"`
		CapacityProviderName string `json:"capacityProviderName"`
		ClusterArn           string `json:"clusterArn"`
		Containers           []struct {
			ContainerArn      string `json:"containerArn"`
			LastStatus        string `json:"lastStatus"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			TaskArn           string `json:"taskArn"`
			NetworkInterfaces []struct {
				AttachmentID       string `json:"attachmentId"`
				PrivateIpv4Address string `json:"privateIpv4Address"`
			} `json:"networkInterfaces"`
			CPU string `json:"cpu"`
		} `json:"containers"`
		CPU           string    `json:"cpu"`
		CreatedAt     time.Time `json:"createdAt"`
		DesiredStatus string    `json:"desiredStatus"`
		Group         string    `json:"group"`
		LaunchType    string    `json:"launchType"`
		LastStatus    string    `json:"lastStatus"`
		Memory        string    `json:"memory"`
		Overrides     struct {
			ContainerOverrides []struct {
				Name string `json:"name"`
			} `json:"containerOverrides"`
		} `json:"overrides"`
		PlatformVersion   string    `json:"platformVersion"`
		StartedBy         string    `json:"startedBy"`
		TaskArn           string    `json:"taskArn"`
		TaskDefinitionArn string    `json:"taskDefinitionArn"`
		UpdatedAt         time.Time `json:"updatedAt"`
		Version           int       `json:"version"`
	} `json:"detail"`
}

var env Env

func init() {
	env = getEnv()
	fmt.Printf("config: %+v\n", env)
}

func HandleLambdaEvent(event MyEvent) error {
	if env.Debug {
		b, _ := json.Marshal(event)
		fmt.Printf("event: %s\n", string(b))
	}

	containerName := event.Detail.Containers[0].Name
	desiredStatus := event.Detail.DesiredStatus
	lastStatus := event.Detail.Containers[0].LastStatus

	for _, ignoreContainerName := range env.IgnoreContainerNames {
		if ignoreContainerName == containerName {
			fmt.Printf("Ignored. container name: %s\n", containerName)
			return nil
		}
	}
	if desiredStatus != lastStatus {
		if lastStatus == "STOPPED" {
			sess := session.Must(session.NewSession())
			svc := ecs.New(sess)

			describeTasksInput := &ecs.DescribeTasksInput{
				Cluster: &event.Detail.ClusterArn,
				Tasks:   []*string{&event.Detail.TaskArn},
			}

			describeTasks, err := svc.DescribeTasks(describeTasksInput)
			if err != nil {
				panic(err)
			}

			slackPayload := generateSlackPayload(event.Detail.Containers[0].Name, desiredStatus, lastStatus, describeTasks)
			errs := sendSlack(slackPayload)
			if errs != nil {
				panic(err)
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
