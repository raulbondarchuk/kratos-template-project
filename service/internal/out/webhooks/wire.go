package webhooks

import (
	webhook "service/internal/out/webhooks/webhook"

	"github.com/google/wire"
)

// ProviderSet is webhooks providers.
var ProviderSet = wire.NewSet(
	// Webhook
	webhook.NewClient,
)
