package main

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strconv"
)

func generateSlackPayload(containerName, desiredStatus, lastStatus string, describeTasks *ecs.DescribeTasksOutput) slack.Payload {
	var reason string
	var exitCode int64

	// Search abnormal container
	for _, container := range describeTasks.Tasks[0].Containers {
		if *container.Name == containerName {
			reason = *container.Reason
			exitCode = *container.ExitCode
		}
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{
		Title: "Container Name",
		Value: containerName,
		Short: false,
	}).AddField(slack.Field{
		Title: "Desired Status",
		Value: desiredStatus,
		Short: true,
	}).AddField(slack.Field{
		Title: "Last Status",
		Value: lastStatus,
		Short: true,
	}).AddField(slack.Field{
		Title: "Reason",
		Value: reason,
		Short: false,
	}).AddField(slack.Field{
		Title: "Exit Code",
		Value: strconv.FormatInt(exitCode, 10),
		Short: false,
	})

	color := "danger"
	attachment.Color = &color

	footer := "Powered by rluisr/ecssc"
	authorLink := "https://github.com/rluisr/ecssc"
	attachment.Footer = &footer
	attachment.AuthorLink = &authorLink

	return slack.Payload{
		Text:        fmt.Sprintf("The container *%s* state is *%s*", containerName, lastStatus),
		Username:    "ecs-state-check",
		Channel:     env.SlackChannelName,
		Attachments: []slack.Attachment{attachment},
	}

}

func sendSlack(payload slack.Payload) []error {
	if env.SlackIconEmoji != "" {
		payload.IconEmoji = env.SlackIconEmoji
	} else {
		payload.IconUrl = env.SlackIconURL
	}

	return slack.Send(env.SlackWebhookURL, "", payload)
}
