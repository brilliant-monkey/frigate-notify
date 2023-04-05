package config

type PushConfig struct {
	Subscriber string      `yaml:"subscriber"`
	VAPID      VapidConfig `yaml:"vapid"`
}
