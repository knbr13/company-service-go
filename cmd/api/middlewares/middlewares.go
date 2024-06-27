package middlewares

import "github.com/knbr13/company-service-go/config"

type Middlewares struct {
	cfg *config.Config
}

func NewMiddlewares(cfg *config.Config) *Middlewares {
	return &Middlewares{cfg: cfg}
}
