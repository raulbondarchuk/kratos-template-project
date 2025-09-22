package broker

import (
	mymqtt "service/pkg/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

type Broker struct {
	log *log.Helper
}

// NewBroker creates a new Broker instance with the given Usecase and logger
func NewBroker(logger log.Logger) *Broker {
	return &Broker{
		log: log.NewHelper(logger),
	}
}

// processMessage processes the message received from the MQTT broker
func (b *Broker) processMessage(client mqtt.Client, message mqtt.Message) {
	// MOCK
	mymqtt.MockMQTT_ProcessMessage(message.Topic(), string(message.Payload()))
	// TODO: Implement the logic to process the message
}
