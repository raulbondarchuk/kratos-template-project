package broker

import (
	"service/internal/conf/v1"
	template_biz "service/internal/feature/template/biz"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"

	mymqtt "service/pkg/mqtt"
	"service/pkg/utils"
)

type Broker struct {
	uc  *template_biz.TemplateUsecase
	log *log.Helper
}

func NewBroker(uc *template_biz.TemplateUsecase, logger log.Logger) *Broker {
	return &Broker{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (b *Broker) StartMQTT(data *conf.Data) {
	server := data.Mqtt.Source
	clientid := data.Mqtt.ClientId
	maxReconnectInterval := data.Mqtt.MaxReconnectInterval
	topics := data.Mqtt.Topics

	username := utils.EnvFirst("MQTT_USERNAME")
	password := utils.EnvFirst("MQTT_PASSWORD")

	b.log.Info("Starting MQTT broker...")
	mymqtt.StartMQTT(server, username, password, &clientid, topics, b.processMessage, maxReconnectInterval.AsDuration())
}

func (b *Broker) processMessage(client mqtt.Client, message mqtt.Message) {

	if err := b.uc.ReceiveTemplate(message.Topic(), string(message.Payload())); err != nil {
		b.log.Errorf("Failed to process message on topic %s: %v", message.Topic(), err)
	}
}
