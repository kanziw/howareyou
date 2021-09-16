package main

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"github.com/kanziw/howareyou/config"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	setting := config.NewSetting()
	_ = slack.New(setting.SlackBotClientSecret)
}
