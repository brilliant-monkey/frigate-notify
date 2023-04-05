package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/brilliant-monkey/frigate-notify/config"
	"github.com/brilliant-monkey/frigate-notify/http"
	"github.com/brilliant-monkey/frigate-notify/types"
	"github.com/brilliant-monkey/go-app"
	"github.com/brilliant-monkey/go-kafka-client"
	"github.com/brilliant-monkey/notify-go"
)

func handleEvent(message []byte) error {
	return nil
}

type NotificationHandler struct {
	pushConfig    *config.PushConfig
	client        *kafka.KafkaClient
	notifier      *notify.Notifier
	subscriptions []*webpush.Subscription
}

func NewNotificationHandler(config config.AppConfig) *NotificationHandler {
	client := kafka.NewKafkaClient(&config.Kafka)
	err := client.TestConnection()
	if err != nil {
		log.Fatalln("Failed to connect to Kafka", err)
	}

	notifier, err := notify.NewNotifier(&config.NotifierConfig)
	if err != nil {
		log.Fatalln("Failed to initialize notifier:", err)
	}

	subscriptions := []*webpush.Subscription{}

	return &NotificationHandler{
		pushConfig:    &config.PushConfig,
		client:        client,
		notifier:      notifier,
		subscriptions: subscriptions,
	}
}

func (handler *NotificationHandler) Subscribe(request []byte) {
	s := &webpush.Subscription{}
	json.Unmarshal(request, s)

}

func (handler *NotificationHandler) Notify(message string) {
	log.Printf("Pushing %s to %v subscribers... ", message, len(handler.subscriptions))
	for _, subscription := range handler.subscriptions {
		resp, err := webpush.SendNotification([]byte(message), subscription, &webpush.Options{
			Subscriber:      fmt.Sprintf("mailto:%s", handler.pushConfig.Subscriber),
			VAPIDPublicKey:  handler.pushConfig.VAPID.PublicKey,
			VAPIDPrivateKey: handler.pushConfig.VAPID.PrivateKey,
			TTL:             30,
		})
		if err != nil {
			log.Println("Failed notify", err)
		}
		defer resp.Body.Close()
	}
}

func (handler *NotificationHandler) Start() error {
	fun := func(message []byte) error {
		var payload types.FrigateEventPayload
		err := json.Unmarshal(message, &payload)
		if err != nil {
			log.Println("Failed to unmarshal event payload.", err)
			return nil
		}

		result, err := handler.notifier.RunTemplate(payload)
		if err != nil {
			log.Println("Failed to run template.", err)
			return nil
		}

		// Notify with result
		handler.Notify(string(result))
		return nil
	}

	return handler.client.Consume(fun)
}

func (handler *NotificationHandler) Stop() error {
	return handler.client.Stop()
}

func main() {
	a := app.NewApp()

	var appConfig config.AppConfig
	a.LoadConfig("CONFIG_PATH", &appConfig)
	handler := NewNotificationHandler(appConfig)

	a.Go(func() error {
		err := handler.Start()
		log.Println("start err", err)
		return err
	})

	apiServer := http.NewAPIServer(&appConfig)
	a.Go(func() error {
		return apiServer.Start()
	})

	a.Start(func() error {
		err := apiServer.Stop()
		if err != nil {
			return err
		}
		return handler.Stop()
	})
}
