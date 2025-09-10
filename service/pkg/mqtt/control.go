package mqtt

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttStarted bool
	mu          sync.Mutex
)

func StartMQTT(server, user, pass string, clientID *string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration) error {
	if topics == nil {
		log.Printf("[MQTT] topics are nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if !mqttStarted {
		mqttStarted = true
		go InitMosquitero(server, user, pass, clientID, topics, handler, maxReconnectInterval)
		return nil
	}
	return fmt.Errorf("[MQTT] already running")
}

func StopMQTT(topics []string) {
	if mqtinstance != nil {
		if err := mqtinstance.Unsubscribe(topics); err != nil {
			fmt.Printf("[MQTT] Unsubscribe error: %v\n", err)
		} else {
			fmt.Println("[MQTT] Unsubscribed successfully")
		}
	}
	mqttStarted = false
}

func (m *Mosquitero) CheckConnection() {
	if !m.client.IsConnected() {
		if token := m.client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] Reconnect failed: %s", token.Error())
		}
	}
}
