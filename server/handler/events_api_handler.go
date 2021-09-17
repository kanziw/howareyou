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
)

const (
	ctxTagsEventDataType     = "evt.data.type"
	ctxTagsKeyInnerEventType = "evt.data.inner_event.type"
	ctxTagsKeyInnerEventData = "evt.data.inner_event.data"
)

var errUnexpectedInnerEventData = errors.New("unexpected evt.data.inner_event.data")

func EventsAPIHandler(ctx context.Context, eventsAPIEvent slackevents.EventsAPIEvent, api *slack.Client) error {
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
			return sendHelpMessage(ctx, api, d.Channel)
		}

		command := strings.ToLower(ss[1])
		args := ss[2:]
		switch command {
		case "start":
			if len(args) == 0 {
				break
			}

			userGroup := args[0]
			_, err := api.GetUserGroupMembersContext(ctx, userGroup)
			if err != nil {
				if err.Error() == "no_such_subteam" {
					_ = sendMessage(ctx, api, d.Channel, userGroup+" is not a user group")
					return sendHelpMessage(ctx, api, d.Channel)
				}
				return errors.WithStack(err)
			}

			// TODO: Upsert Schedule into DB

			if err := sendMessage(ctx, api, d.Channel, "Let's start HowAreYou!"); err != nil {
				return errors.WithStack(err)
			}
		}
		return sendHelpMessage(ctx, api, d.Channel)
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

func sendMessage(ctx context.Context, api *slack.Client, channel, msg string) error {
	if _, _, _, err := api.SendMessageContext(ctx, channel, slack.MsgOptionText(msg, false)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// TODO
func sendHelpMessage(ctx context.Context, api *slack.Client, channel string) error {
	if err := sendMessage(ctx, api, channel, "help message"); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
