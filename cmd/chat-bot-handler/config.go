package main

// Config application configuration
type Config struct {
	LogLevel   string `long:"log-level" description:"Log level: panic, fatal, warn or warning, info, debug" env:"LOG_LEVEL" required:"true"`
	LogJSON    bool   `long:"log-json" description:"Enable force log format JSON" env:"LOG_JSON"`
	HttpListen string `long:"http-listen" description:"HTTP listen port" env:"HTTP_LISTEN"`

	TgApiKey string `long:"tg-api-key" description:"Telegram api key" env:"TG_API_KEY"`

	EnablePprof bool `long:"enable-pprof" description:"Enable pprof server" env:"ENABLE_PPROF"`
}
