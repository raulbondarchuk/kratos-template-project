package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Mosquitero struct {
	client           mqtt.Client
	subscribedTopics []string
	handler          mqtt.MessageHandler
}

var mqtinstance *Mosquitero

func GetMosquitero() *Mosquitero {
	return mqtinstance
}

func (m *Mosquitero) GetClient() mqtt.Client {
	return m.client
}
