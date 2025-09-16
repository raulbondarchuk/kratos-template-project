package broker

import mqtt "github.com/eclipse/paho.mqtt.golang"

func (b *Broker) processMessage(client mqtt.Client, message mqtt.Message) {

	if err := b.uc.ReceiveTemplate(message.Topic(), string(message.Payload())); err != nil {
		b.log.Errorf("Failed to process message on topic %s: %v", message.Topic(), err)
	}
}
