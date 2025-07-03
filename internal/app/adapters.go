package app

import (
	httpAlerter "github.com/dsaime/auth-api/internal/adapter/http_alerter"
	"github.com/dsaime/auth-api/internal/service"
)

type adapters struct {
	alerter service.Alerter
}

func initAdapters(cfg Config) *adapters {

	return &adapters{
		alerter: &httpAlerter.HttpAlerter{
			Webhook: cfg.Webhook,
		},
	}
}
