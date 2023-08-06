package main

type Config struct {
	LogLevel   string `long:"log-level" description:"Log level: panic, fatal, warn or warning, info, debug" env:"LOG_LEVEL" required:"true"`
	LogJSON    bool   `long:"log-json" description:"Enable force log format JSON" env:"LOG_JSON"`
	HttpListen string `long:"http-listen" description:"HTTP listen port" env:"HTTP_LISTEN"`

	KafkaBroker          string `long:"kafka-broker" description:"Kafka broker" env:"KAFKA_BROKER"`
	KafkaMaxMessageBytes int    `long:"kafka-max-size-message" description:"Max size message for Kafka" env:"KAFKA_MAX_MESSAGE_BYTES" required:"true"`
	KafkaMaxRetry        int    `long:"kafka-max-retry" description:"Max retry count to connect to Kafka" env:"KAFKA_MAX_RETRY" required:"true"`

	TgApiKey string `long:"tg-api-key" description:"Telegram api key" env:"TG_API_KEY"`

	EnablePprof bool `long:"enable-pprof" description:"Enable pprof server" env:"ENABLE_PPROF"`
}
