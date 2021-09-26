package config

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Setting struct {
	IsDebug       bool
	SlackBotToken string
	SlackRTMToken string

	DBHost            string
	DBPort            int
	DBName            string
	DBUser            string
	DBPassword        string
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
}

func NewSetting() Setting {
	return Setting{
		IsDebug:       mustParseBool("IS_DEBUG", false),
		SlackBotToken: getEnv("SLACK_BOT_TOKEN", ""),
		SlackRTMToken: getEnv("SLACK_RTM_TOKEN", ""),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            mustParseInt("DB_PORT", 3306),
		DBName:            getEnv("DB_NAME", "howareyou"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", "root"),
		DBMaxIdleConns:    mustParseInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns:    mustParseInt("DB_MAX_OPEN_CONNS", 10),
		DBConnMaxLifetime: mustParseDuration("DB_CONN_MAX_LIFETIME_DURATION", "3m"),
	}
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultValue
	}
	if v == "" {
		log.Fatalf("env %s is not set", key)
	}
	return v
}

func mustParseBool(key string, defaultValue bool) bool {
	v := getEnv(key, strconv.FormatBool(defaultValue))
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.WithError(err).Fatalf("unexpected bool value %s with key %s", v, key)
	}
	return b
}

func mustParseInt(key string, defaultValue int) int {
	v := getEnv(key, strconv.Itoa(defaultValue))
	i, err := strconv.Atoi(v)
	if err != nil {
		log.WithError(err).Fatalf("unexpected int value %s with key %s", v, key)
	}
	return i
}

func mustParseDuration(key, defaultValue string) time.Duration {
	v := getEnv(key, defaultValue)
	d, err := time.ParseDuration(v)
	if err != nil {
		log.WithError(err).Fatalf("unexpected int value %s with key %s", v, key)
	}
	return d
}
