package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
)

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

var config Env

func init() {
	config = getConfig()
	fmt.Printf("config: %+v\n", config)
}

func HandleLambdaEvent(event MyEvent) error {
	containerName := event.Detail.Containers[0].Name
	desiredStatus := event.Detail.DesiredStatus
	lastStatus := event.Detail.Containers[0].LastStatus

	for _, ignoreContainerName := range config.IgnoreContainerNames {
		if ignoreContainerName == containerName {
			fmt.Printf("Ignored. container name: %s\n", containerName)
			return nil
		}
	}
	if desiredStatus != lastStatus {
		if lastStatus == "STOPPED" {
			sendSlack(event.Detail.Containers[0].Name, desiredStatus, lastStatus)
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
