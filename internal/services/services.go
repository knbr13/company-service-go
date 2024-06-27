package services

import (
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/knbr13/company-service-go/internal/repositories"
)

type Services struct {
	Users     *UserService
	Companies *CompanyService
}

func NewServices(db *sql.DB, producer sarama.SyncProducer, errCh chan<- error) *Services {
	repos := repositories.NewRepositories(db)

	return &Services{
		Users:     &UserService{Repos: &repos},
		Companies: &CompanyService{Repos: &repos, Producer: producer, ErrCh: errCh},
	}
}
