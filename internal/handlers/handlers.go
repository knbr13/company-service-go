package handlers

import (
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/knbr13/company-service-go/config"
	"github.com/knbr13/company-service-go/internal/services"
)

type Handlers struct {
	Users     *UserHandler
	Companies *CompanyHandler
}

func NewHandlers(db *sql.DB, cfg *config.Config, producer sarama.SyncProducer, errCh chan<- error) *Handlers {
	srvcs := services.NewServices(db, producer, errCh)

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
