package main

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/kanziw/howareyou/config"
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

	go func() {
		for evt := range client.Events {
			log := logrus.StandardLogger().WithFields(logrus.Fields{
				"evt.type": evt.Type,
			})
			switch evt.Type {
			case socketmode.RequestTypeHello:
			case socketmode.EventTypeConnecting:
				log.Debug("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Debug("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Debug("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.WithField("evt", evt).Error("unknown event type:" + evt.Type)
					continue
				}
				client.Ack(*evt.Request)

				log = log.
					WithFields(logrus.Fields{
						"evt.data": logrus.Fields{
							"type": eventsAPIEvent.Type,
							"inner_event": logrus.Fields{
								"type": eventsAPIEvent.InnerEvent.Type,
							},
						},
					})

				switch eventsAPIEvent.InnerEvent.Type {
				case slackevents.ReactionAdded:
					d, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
					if !ok {
						log.WithField("evt.data.inner_event.data", d).Error("unexpected evt.data.inner_event.data")
						continue
					}
					log.Infof("%s User react using %s on %s's message", d.User, d.Reaction, d.ItemUser)
				default:
					log.Info("unsupported Events API event received")
					client.Debugf("unsupported Events API event received")
				}
			// TODO
			case socketmode.EventTypeInteractive:
			case socketmode.EventTypeSlashCommand:
			default:
				log.Info("unexpected event type received")
			}
		}
	}()

	if err := client.Run(); err != nil {
		logrus.Fatal(err)
	}
}
