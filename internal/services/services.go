package services

import (
	"database/sql"

	"github.com/knbr13/company-service-go/internal/repositories"
)

type Services struct {
	Users     *UserService
	Companies *CompanyService
}

func NewServices(db *sql.DB) *Services {
	repos := repositories.NewRepositories(db)

	return &Services{
		Users:     &UserService{Repos: &repos},
		Companies: &CompanyService{Repos: &repos},
	}
}
