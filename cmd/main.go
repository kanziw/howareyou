package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/kanziw/go-slack"
	"github.com/kanziw/howareyou/config"
	"github.com/kanziw/howareyou/mysql"
	"github.com/kanziw/howareyou/service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()

	setting := config.NewSetting()

	db, err := mysql.GetDB(setting)
	if err != nil {
		log.WithError(err).Fatal("mysql.GetDB")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(err)
		}
	}()

	s := slack.NewSocketServer(
		setting.SlackBotToken,
		setting.SlackRTMToken,
		slack.WithDebug(false),
	)
	svc := service.New(db, s.SlackAPI())

	s.OnAppMentionCommand("start", func(ctx context.Context, d *slack.AppMentionEvent, api *slack.Client, args []string) error {
		if len(args) == 0 {
			return errInvalidCommand(d.Channel)
		}

		userGroup := args[0]
		// for test in free plan
		//if _, err := api.GetUserGroupMembersContext(ctx, userGroup); err != nil {
		//	if err.Error() == "no_such_subteam" {
		//		// It's not important. Ignore
		//		_ = slack.SendMessage(ctx, api, d.Channel, userGroup+" is not a user group")
		//		return errInvalidCommand(d.Channel)
		//	}
		//	return errors.WithStack(err)
		//}
		return svc.StartHowAreYou(ctx, d.Channel, userGroup)
	})

	go s.Listen()

	if err := s.Run(); err != nil {
		logrus.Fatal(err)
	}
}

func errInvalidCommand(channel string) error {
	return slack.NewSlackError(errors.WithStack(slack.ErrInvalidCommand), slack.WithChannel(channel))
}
