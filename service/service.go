package service

import (
	"context"

	"github.com/slack-go/slack"
)

type Service interface {
	StartHowAreYou(ctx context.Context, channel, userGroup string) error
}

type DefaultService struct {
	api *slack.Client
}

func (s *DefaultService) StartHowAreYou(ctx context.Context, channel, userGroup string) error {
	return nil
}

func New(api *slack.Client) Service {
	return &DefaultService{
		api: api,
	}
}
