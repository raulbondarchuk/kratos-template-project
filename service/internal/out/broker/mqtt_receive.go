package broker

import (
	mymqtt "service/pkg/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// processMessage processes the message received from the MQTT broker
func (b *Broker) processMessage(client mqtt.Client, message mqtt.Message) {
	// MOCK
	mymqtt.MockMQTT_ProcessMessage(message.Topic(), string(message.Payload()))
	// TODO: Implement the logic to process the message
}
