package main

import (
	"io/ioutil"
	l "log"
	"os"
	"strings"
)

func init() {
	discordConf := &discordConfig{
		BotToken: configOption{
			EnvironmentKey: "DMSB_DISCORT_BOT_TOKEN",
			SecretFile:     "discord-bot-token",
			FlagKey:        "botToken",
		}.Resolve(),
	}

	loggerConf := &loggerConfig{
		Level: configOption{
			EnvironmentKey: "DMSB_LOGGER_LOG_LEVEL",
			FlagKey:        "logLevel",
			DefaultValue:   "debug",
		}.Resolve(),
		File: configOption{
			EnvironmentKey: "DMSB_LOGGER_LOG_FILE",
			FlagKey:        "logFile",
		}.Resolve(),
	}

	storageConf := &storageConfig{
		S3StorageBucketURI: configOption{
			EnvironmentKey: "DMSB_STORAGE_S3_STORAGE_BUCKET_URI",
			FlagKey:        "s3StorageBucketURI",
		}.Resolve(),
	}

	dbConf := &dbConfig{
		RootPassword: configOption{
			EnvironmentKey: "DMSB_DB_MYSQL_ROOT_PASSWORD",
			SecretFile:     "mysql-root-password",
			FlagKey:        "mySQLRootPassword",
		}.Resolve(),
	}

	conf = &config{
		Discord: discordConf,
		Logger:  loggerConf,
		Storage: storageConf,
		DB:      dbConf,
	}
}

var conf *config

type config struct {
	Discord *discordConfig
	Logger  *loggerConfig
	Storage *storageConfig
	DB      *dbConfig
}

type discordConfig struct {
	BotToken string
}

type loggerConfig struct {
	Level string
	File  string
}

type dbConfig struct {
	RootPassword string
}

type storageConfig struct {
	S3StorageBucketURI string
}

type configOption struct {
	EnvironmentKey string
	SecretFile     string
	FlagKey        string
	DefaultValue   string
	UsageText      string
	Required       bool
}

func (o configOption) Resolve() string {
	val := ""
	// val = *flag.String(o.FlagKey, "", o.UsageText)
	if len(val) == 0 && len(o.EnvironmentKey) != 0 {
		val = os.Getenv(o.EnvironmentKey)
	}
	if len(val) == 0 && len(o.SecretFile) != 0 {
		b, err := ioutil.ReadFile("secrets/" + o.SecretFile)
		if err == nil {
			val = strings.TrimSpace(string(b))
		}
	}
	if len(val) == 0 && o.Required {
		l.Fatalf("config variable `%s` is required but undefined", o.FlagKey)
	}
	if len(val) == 0 {
		val = o.DefaultValue
	}
	return val
}
