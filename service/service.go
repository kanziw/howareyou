package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/kanziw/go-slack"
)

type Service interface {
	StartHowAreYou(ctx context.Context, channel, userGroup string) error
}

type DefaultService struct {
	api *slack.Client
}

func (s *DefaultService) StartHowAreYou(ctx context.Context, channel, userGroup string) error {
	// TODO: Upsert Schedule into DB

	if err := slack.SendMessage(ctx, s.api, channel, "Let's start HowAreYou with "+userGroup); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *DefaultService) SendHowAreYou(ctx context.Context, channel, userGroup string) error {
	if err := slack.SendMessage(
		ctx,
		s.api,
		channel,
		fmt.Sprintf("Hi %s :) How are you today?", userGroup),
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func New(api *slack.Client) Service {
	return &DefaultService{
		api: api,
	}
}
