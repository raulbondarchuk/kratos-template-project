package mqtt

func MockMQTT_ProcessMessage(topic, message string) {
	mqttLogger.Infof("MockMQTT_ProcessMessage: %s, %s", topic, message)
}
