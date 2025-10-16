package webhook

import (
	"service/internal/conf/v1"
	"time"

	"github.com/go-resty/resty/v2"
)

type clientImpl struct {
	client  *resty.Client
	baseURL string

	loginRoute string
	userRoute  string
}

type Client interface {
}

func NewClient(cfg *conf.Webhooks) (Client, error) {
	timeout := cfg.Webhook.Timeout.AsDuration()
	if timeout <= 5*time.Second {
		timeout = 30 * time.Second
	}

	cli := resty.New().SetTimeout(timeout)

	impl := &clientImpl{
		client:     cli,
		baseURL:    cfg.Webhook.Url,
		loginRoute: cfg.Webhook.Routes.Route1,
		userRoute:  cfg.Webhook.Routes.Route2,
	}

	return impl, nil
}
