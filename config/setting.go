package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Setting struct {
	SlackBotClientSecret string
}

func NewSetting() Setting {
	return Setting{
		SlackBotClientSecret: getEnv("SLACK_BOT_CLIENT_SECRET", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultValue
	}
	if v == "" {
		log.Fatalf("env %s is not set", key)
	}
	return v
}
