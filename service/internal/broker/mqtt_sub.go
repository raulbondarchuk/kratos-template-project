package broker

import (
	"service/internal/conf/v1"

	mymqtt "service/pkg/mqtt"
	"service/pkg/utils"
)

// StartMQTT starts the MQTT broker
func (b *Broker) StartMQTT(data *conf.Data) {

	if data.Mqtt == nil {
		b.log.Info("[MQTT] [SKIPPED] Broker is not configured, skipping...")
		return
	} else if !data.Mqtt.Active {
		b.log.Info("[MQTT] [SKIPPED] Broker is inactive, skipping...")
		return
	}

	server := data.Mqtt.Source
	clientid := data.Mqtt.ClientId
	maxReconnectInterval := data.Mqtt.MaxReconnectInterval
	topics := data.Mqtt.Topics

	username := utils.EnvFirst("MQTT_USERNAME")
	password := utils.EnvFirst("MQTT_PASSWORD")

	b.log.Info("Starting MQTT broker...")
	mymqtt.StartMQTT(server, username, password, &clientid, topics, b.processMessage, maxReconnectInterval.AsDuration(), b.log.Logger())
}
