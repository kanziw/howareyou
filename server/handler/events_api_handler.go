package handler

import (
	"context"
	"fmt"
	"strings"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/kanziw/howareyou/service"
)

const (
	ctxTagsEventDataType     = "evt.data.type"
	ctxTagsKeyInnerEventType = "evt.data.inner_event.type"
	ctxTagsKeyInnerEventData = "evt.data.inner_event.data"
)

var errUnexpectedInnerEventData = errors.New("unexpected evt.data.inner_event.data")

func EventsAPIHandler(
	ctx context.Context,
	eventsAPIEvent slackevents.EventsAPIEvent,
	api *slack.Client,
	svc service.Service,
) error {
	tags := grpc_ctxtags.Extract(ctx)
	tags.Set(ctxTagsEventDataType, eventsAPIEvent.Type)
	tags.Set(ctxTagsKeyInnerEventType, eventsAPIEvent.InnerEvent.Type)

	switch eventsAPIEvent.InnerEvent.Type {
	case slackevents.AppMention:
		d, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
		if !ok {
			tags.Set(ctxTagsKeyInnerEventData, d)
			return errors.WithStack(errUnexpectedInnerEventData)
		}
		tags.Set(ctxTagsKeyInnerEventData, logrus.Fields{
			"user":        d.User,
			"channel":     d.Channel,
			"text":        d.Text,
			"description": fmt.Sprintf("%s User mention in channel %s with text %s", d.User, d.Channel, d.Text),
		})

		ss := strings.Split(strings.TrimSpace(d.Text), " ")
		if len(ss) < 2 {
			return service.SendHelpMessage(ctx, api, d.Channel)
		}

		command := strings.ToLower(ss[1])
		args := ss[2:]
		switch command {
		case "start":
			if len(args) == 0 {
				break
			}

			userGroup := args[0]
			if _, err := api.GetUserGroupMembersContext(ctx, userGroup); err != nil {
				if err.Error() == "no_such_subteam" {
					// It's not important. Ignore
					_ = service.SendMessage(ctx, api, d.Channel, userGroup+" is not a user group")
					return service.SendHelpMessage(ctx, api, d.Channel)
				}
				return errors.WithStack(err)
			}
			return svc.StartHowAreYou(ctx, d.Channel, userGroup)
		}
		return service.SendHelpMessage(ctx, api, d.Channel)
	case slackevents.ReactionAdded:
		d, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
		if !ok {
			tags.Set(ctxTagsKeyInnerEventData, d)
			return errors.WithStack(errUnexpectedInnerEventData)
		}
		tags.Set(ctxTagsKeyInnerEventData, logrus.Fields{
			"user":        d.User,
			"reaction":    d.Reaction,
			"item_user":   d.ItemUser,
			"description": fmt.Sprintf("%s User react using %s on %s's message", d.User, d.Reaction, d.ItemUser),
		})
		return nil
	}

	return errors.New("unsupported Events API event received")
}
