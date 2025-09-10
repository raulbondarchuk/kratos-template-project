package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (m *Mosquitero) Subscribe(topics []string, handler mqtt.MessageHandler) {
	m.subscribedTopics = topics
	m.handler = handler
	for _, topic := range topics {
		if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			log.Printf("❌ [MQTT] Subscribe error %s: %s", topic, token.Error())
		} else {
			log.Printf("✅ [MQTT] Subscribed: %s", topic)
		}
	}
}

func (m *Mosquitero) Unsubscribe(topics []string) error {
	token := m.client.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}
