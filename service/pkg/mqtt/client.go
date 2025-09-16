package mqtt

import (
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

var mqtonce sync.Once

func InitMosquitero(server, username, password string, clientid *string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration, logger log.Logger) *Mosquitero {
	h := log.NewHelper(logger)

	if clientid == nil {
		id := uuid.New().String()[:6]
		clientid = &id
		h.Infof("✅ [MQTT] Random Client ID: %s", *clientid)
	} else {
		*clientid += "_" + uuid.New().String()[:6]
		h.Infof("✅ [MQTT] Client ID: %s", *clientid)
	}
	return internalInitMosquitero(server, username, password, *clientid, topics, handler, maxReconnectInterval, logger)
}

func internalInitMosquitero(server, username, password, clientid string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration, logger log.Logger) *Mosquitero {
	h := log.NewHelper(logger)

	mqtonce.Do(func() {
		opts := mqtt.NewClientOptions().
			AddBroker(server).
			SetClientID(clientid).
			SetUsername(username).
			SetPassword(password).
			SetAutoReconnect(true).
			SetMaxReconnectInterval(maxReconnectInterval)

		opts.OnConnect = func(c mqtt.Client) {
			h.Infof("✅ [MQTT] Connected to broker")
			if mqtinstance != nil && len(mqtinstance.subscribedTopics) > 0 && mqtinstance.handler != nil {
				mqtinstance.Subscribe(mqtinstance.subscribedTopics, mqtinstance.handler)
			} else {
				mqtinstance.Subscribe(topics, handler)
			}
		}
		opts.OnConnectionLost = func(c mqtt.Client, err error) {
			h.Infof("[MQTT] Connection lost: %v", err)
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			h.Fatal("[MQTT] Connection error: " + token.Error().Error())
		}

		mqtinstance = &Mosquitero{client: client}
	})
	return mqtinstance
}
