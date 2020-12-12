package main

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
)

func sendSlack(containerName, desiredStatus, lastStatus string) {
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
	})

	color := "danger"
	attachment.Color = &color

	footer := "Powered by ECSSC"
	authorLink := "https://github.com/rluisr/ecssc"
	attachment.Footer = &footer
	attachment.AuthorLink = &authorLink

	payload := slack.Payload{
		Text:        fmt.Sprintf("The container *%s* state is *%s*", containerName, lastStatus),
		Username:    "ecs-state-check",
		Channel:     config.SlackChannelName,
		Attachments: []slack.Attachment{attachment},
	}

	if config.SlackIconEmoji != "" {
		payload.IconEmoji = config.SlackIconEmoji
	} else {
		payload.IconUrl = config.SlackIconURL
	}

	errs := slack.Send(config.SlackWebhookURL, "", payload)
	if len(errs) > 0 {
		panic(errs)
	}
}
