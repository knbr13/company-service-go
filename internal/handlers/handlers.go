package handlers

import (
	"database/sql"

	"github.com/knbr13/company-service-go/config"
	"github.com/knbr13/company-service-go/internal/services"
)

type Handlers struct {
	Users     *UserHandler
	Companies *CompanyHandler
}

func NewHandlers(db *sql.DB, cfg *config.Config) *Handlers {
	srvcs := services.NewServices(db)

	return &Handlers{
		Users: &UserHandler{
			Services: srvcs,
			Cfg:      cfg,
		},
		Companies: &CompanyHandler{
			Services: srvcs,
		},
	}
}
