package service

import (
	"context"
	"fmt"

	"github.com/kanziw/go-slack"

	"github.com/pkg/errors"
)

type Service interface {
	StartHowAreYou(ctx context.Context, channel, userGroup string) error
}

type DefaultService struct {
	api *slack.Client
}

func (s *DefaultService) StartHowAreYou(ctx context.Context, channel, userGroup string) error {
	// TODO: Upsert Schedule into DB

	if err := slack.SendMessage(ctx, s.api, channel, "Let's start HowAreYou!"); err != nil {
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
