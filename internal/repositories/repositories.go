package repositories

import "database/sql"

type Repositories struct {
	User    UserRepository
	Company CompanyRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		User:    UserRepository{DB: db},
		Company: CompanyRepository{DB: db},
	}
}
