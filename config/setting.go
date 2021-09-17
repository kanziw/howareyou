package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Setting struct {
	IsDebug       bool
	SlackBotToken string
	SlackRTMToken string
}

func NewSetting() Setting {
	return Setting{
		IsDebug:       parseBool("IS_DEBUG", false),
		SlackBotToken: getEnv("SLACK_BOT_TOKEN", ""),
		SlackRTMToken: getEnv("SLACK_RTM_TOKEN", ""),
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

func parseBool(key string, defaultValue bool) bool {
	v := getEnv(key, strconv.FormatBool(defaultValue))
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.WithError(err).Fatalf("unexpected bool value %s with key %s", v, key)
	}
	return b
}
