package broker

import (
	template_biz "service/internal/feature/template/v1/biz"
	mymqtt "service/pkg/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

type Broker struct {
	uc  *template_biz.TemplateUsecase // <---- Usecase for processing messages logic
	log *log.Helper
}

// NewBroker creates a new Broker instance with the given Usecase and logger
func NewBroker(uc *template_biz.TemplateUsecase, logger log.Logger) *Broker {
	return &Broker{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// processMessage processes the message received from the MQTT broker
func (b *Broker) processMessage(client mqtt.Client, message mqtt.Message) {
	mymqtt.MockMQTT_ProcessMessage(message.Topic(), string(message.Payload()))
}
