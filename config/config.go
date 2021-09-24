package config

import (
	"github.com/kanziw/go-slack"
	"github.com/kanziw/howareyou/service"
)

type Config struct {
	setting Setting

	api    *slack.Client
	client *slack.SocketClient

	svc service.Service
}

func (c *Config) SlackAPI() *slack.Client {
	return c.api
}

func (c *Config) SocketClient() *slack.SocketClient {
	return c.client
}

func (c *Config) Service() service.Service {
	return c.svc
}

func New(setting Setting, api *slack.Client, client *slack.SocketClient, svc service.Service) *Config {
	return &Config{
		setting: setting,
		api:     api,
		client:  client,
		svc:     svc,
	}
}
