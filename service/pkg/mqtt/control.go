package mqtt

import (
	"fmt"

	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	mqttStarted bool
	mu          sync.Mutex

	mqttLogger *log.Helper
)

func StartMQTT(server, user, pass string, clientID *string, topics []string, handler mqtt.MessageHandler, maxReconnectInterval time.Duration, logger log.Logger) error {

	if mqttLogger == nil {
		mqttLogger = log.NewHelper(logger)
	}

	if topics == nil {
		mqttLogger.Infof("[MQTT] topics are nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if !mqttStarted {
		mqttStarted = true
		go InitMosquitero(server, user, pass, clientID, topics, handler, maxReconnectInterval, logger)
		return nil
	}
	return fmt.Errorf("[MQTT] already running")
}

func StopMQTT(topics []string) {

	if mqtinstance != nil {
		if err := mqtinstance.Unsubscribe(topics); err != nil {
			mqttLogger.Infof("[MQTT] Unsubscribe error: %v\n", err)
		} else {
			mqttLogger.Infof("[MQTT] Unsubscribed successfully")
		}
	}
	mqttStarted = false
}

func (m *Mosquitero) CheckConnection() {
	if !m.client.IsConnected() {
		if token := m.client.Connect(); token.Wait() && token.Error() != nil {
			mqttLogger.Infof("[MQTT] Reconnect failed: %s", token.Error())
		}
	}
}
