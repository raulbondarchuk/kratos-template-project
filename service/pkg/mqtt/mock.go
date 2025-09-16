package mqtt

func MockMQTT_ProcessMessage(topic, message string) {
	truncatedMessage := message
	if len(message) > 15 {
		truncatedMessage = message[:20] + "..."
	}
	mqttLogger.Infof("[MQTT] [MOCK] received %s: %s", topic, truncatedMessage)
}
