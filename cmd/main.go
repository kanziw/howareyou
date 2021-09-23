package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"

	"github.com/kanziw/go-slack"
	"github.com/kanziw/go-slack/handler"
	"github.com/kanziw/howareyou/config"
	"github.com/kanziw/howareyou/service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	setting := config.NewSetting()

	s := slack.NewSocketServer(
		setting.SlackBotToken,
		setting.SlackRTMToken,
		slack.WithDebug(false),
	)
	svc := service.New(s.SlackAPI())

	s.OnAppMentionCommand("start", func(ctx context.Context, d *slackevents.AppMentionEvent, api *slack.Client, args []string) error {
		if len(args) == 0 {
			return errors.WithStack(handler.ErrInvalidCommand)
		}

		userGroup := args[0]
		if _, err := api.GetUserGroupMembersContext(ctx, userGroup); err != nil {
			if err.Error() == "no_such_subteam" {
				// It's not important. Ignore
				_ = service.SendMessage(ctx, api, d.Channel, userGroup+" is not a user group")
				return errors.WithStack(handler.ErrInvalidCommand)
			}
			return errors.WithStack(err)
		}
		return svc.StartHowAreYou(ctx, d.Channel, userGroup)
	})

	go s.Listen()

	if err := s.Run(); err != nil {
		logrus.Fatal(err)
	}
}
