package main

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"github.com/kanziw/howareyou/config"
	"github.com/kanziw/howareyou/server"
	"github.com/kanziw/howareyou/service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	setting := config.NewSetting()
	api := slack.New(
		setting.SlackBotToken,
		slack.OptionDebug(setting.IsDebug),
		slack.OptionAppLevelToken(setting.SlackRTMToken),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(setting.IsDebug),
	)
	cfg := config.New(
		setting,
		api,
		client,
		service.New(api),
	)

	s := server.NewSocketServer(cfg)
	go s.Listen()

	if err := client.Run(); err != nil {
		logrus.Fatal(err)
	}
}
