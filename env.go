package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Debug                bool     `default:"false"`
	IgnoreContainerNames []string `split_words:"true"`
	SlackWebhookURL      string   `required:"true" split_words:"true"`
	SlackChannelName     string   `required:"true" split_words:"true"`
	SlackUserName        string   `default:"ecs-state-check (ecssc)" split_words:"true"`
	SlackIconURL         string   `default:"https://f.easyuploader.app/eu-prd/upload/20201213004929_72414e62464d46756a47.png" split_words:"true"`
	SlackIconEmoji       string   `split_words:"true"`
}

func getEnv() Env {
	var e Env
	err := envconfig.Process("ecssc", &e)
	if err != nil {
		panic(err)
	}

	return e
}
