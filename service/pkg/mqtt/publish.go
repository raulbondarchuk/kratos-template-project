package mqtt

import (
	"encoding/json"
)

func (m *Mosquitero) Send(topic, payload string) {
	go m.send(topic, 0, payload)
}

func (m *Mosquitero) SendQos(topic string, qos byte, payload string) {
	go m.send(topic, qos, payload)
}

func (m *Mosquitero) send(topic string, qos byte, payload string) {
	token := m.client.Publish(topic, qos, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		mqttLogger.Infof("[MQTT] Error publishing to %s: %s", topic, err)
	}
}

func (m *Mosquitero) SendJSON(topic string, v any) error {
	return m.SendJSONEx(topic, 0, false, v)
}

func (m *Mosquitero) SendJSONEx(topic string, qos byte, retained bool, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	token := m.client.Publish(topic, qos, retained, data)
	token.Wait()
	return token.Error()
}
