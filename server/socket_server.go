package server

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/kanziw/howareyou/config"
	"github.com/kanziw/howareyou/server/handler"
)

type SocketServer interface {
	Listen()
}

type DefaultSocketServer struct {
	cfg *config.Config
}

func (s *DefaultSocketServer) Listen() {
	for evt := range s.cfg.SocketClient().Events {
		ctx := ctxlogrus.ToContext(
			grpc_ctxtags.SetInContext(context.Background(), grpc_ctxtags.NewTags()),
			logrus.WithField("evt.type", evt.Type),
		)
		err := func() error {
			switch evt.Type {
			case socketmode.RequestTypeHello:
			case socketmode.EventTypeConnecting:
				s.cfg.SocketClient().Debugln("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				s.cfg.SocketClient().Debugln("debug", "Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				s.cfg.SocketClient().Debugln("debug", "Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					return errors.New("unknown event type:" + string(evt.Type))
				}
				s.cfg.SocketClient().Ack(*evt.Request)
				if err := handler.EventsAPIHandler(ctx, eventsAPIEvent, s.cfg.SlackAPI()); err != nil {
					s.cfg.SocketClient().Debugf(err.Error())
					return err
				}
			// TODO
			case socketmode.EventTypeInteractive:
			case socketmode.EventTypeSlashCommand:
			default:
				return errors.New("unexpected event type received")
			}
			return nil
		}()
		entry := ctxlogrus.Extract(ctx).WithContext(ctx)
		if err != nil {
			entry.WithField("error", err).Error(err.Error())
			continue
		}
		entry.Info("succeeded")
	}
}

func NewSocketServer(cfg *config.Config) SocketServer {
	return &DefaultSocketServer{cfg: cfg}
}
