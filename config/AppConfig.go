package config

import "github.com/brilliant-monkey/notify-go/config"

type AppConfig struct {
	Kafka          KafkaConfig           `yaml:"kafka"`
	PushConfig     PushConfig            `yaml:"push"`
	NotifierConfig config.NotifierConfig `yaml:"notifier"`
}
