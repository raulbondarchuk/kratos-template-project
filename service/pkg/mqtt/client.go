package mqtt

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var mqtonce sync.Once

func InitMosquitero(server, username, password string, clientid *string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration) *Mosquitero {
	if clientid == nil {
		id := uuid.New().String()[:6]
		clientid = &id
		log.Printf("✅ [MQTT] Random Client ID: %s", *clientid)
	} else {
		*clientid += "_" + uuid.New().String()[:6]
		log.Printf("✅ [MQTT] Client ID: %s", *clientid)
	}
	return internalInitMosquitero(server, username, password, *clientid, topics, handler, maxReconnectInterval)
}

func internalInitMosquitero(server, username, password, clientid string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration) *Mosquitero {
	mqtonce.Do(func() {
		opts := mqtt.NewClientOptions().
			AddBroker(server).
			SetClientID(clientid).
			SetUsername(username).
			SetPassword(password).
			SetAutoReconnect(true).
			SetMaxReconnectInterval(maxReconnectInterval)

		opts.OnConnect = func(c mqtt.Client) {
			log.Println("✅ [MQTT] Connected to broker")
			if mqtinstance != nil && len(mqtinstance.subscribedTopics) > 0 && mqtinstance.handler != nil {
				mqtinstance.Subscribe(mqtinstance.subscribedTopics, mqtinstance.handler)
			} else {
				mqtinstance.Subscribe(topics, handler)
			}
		}
		opts.OnConnectionLost = func(c mqtt.Client, err error) {
			log.Printf("[MQTT] Connection lost: %v", err)
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal("[MQTT] Connection error: " + token.Error().Error())
		}

		mqtinstance = &Mosquitero{client: client}
	})
	return mqtinstance
}
