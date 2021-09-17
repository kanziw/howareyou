package handler

import (
	"context"
	"fmt"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

func EventsAPIHandler(ctx context.Context, eventsAPIEvent slackevents.EventsAPIEvent) error {
	tags := grpc_ctxtags.Extract(ctx)
	tags.Set("evt.data.type", eventsAPIEvent.Type)
	tags.Set("evt.data.inner_event.type", eventsAPIEvent.InnerEvent.Type)

	switch eventsAPIEvent.InnerEvent.Type {
	case slackevents.ReactionAdded:
		d, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.ReactionAddedEvent)
		if !ok {
			tags.Set("evt.data.inner_event.data", d)
			return errors.New("unexpected evt.data.inner_event.data")
		}
		tags.Set("evt.data.inner_event.data", logrus.Fields{
			"user":        d.User,
			"reaction":    d.Reaction,
			"item_user":   d.ItemUser,
			"description": fmt.Sprintf("%s User react using %s on %s's message", d.User, d.Reaction, d.ItemUser),
		})
		return nil
	}

	return errors.New("unsupported Events API event received")
}
