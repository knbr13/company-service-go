package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrCompanyNameAlreadyExists = errors.New("company name already exists")
)

type Company struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	AmountOfEmployees *int   `json:"amount_of_employees"`
	Registered        *bool  `json:"registered"`
}

type CompanyRepository struct {
	DB *sql.DB
}

func (c *CompanyRepository) Insert(ctx context.Context, comp *Company) error {
	_, err := c.DB.Exec(
		"INSERT INTO companies (id, name, description, type, amount_of_employees, registered) VALUES (?,?,?,?,?,?)",
		comp.ID,
		comp.Name,
		comp.Description,
		comp.Type,
		comp.AmountOfEmployees,
		comp.Registered,
	)
	if err != nil {
		if errors.Is(err, &mysql.MySQLError{Number: 1062}) {
			return ErrCompanyNameAlreadyExists
		}
		return err
	}
	return nil
}

func (c *CompanyRepository) Update(ctx context.Context, comp *Company) error {
	_, err := c.DB.Exec(
		"UPDATE companies SET name =?, description =?, type =?, amount_of_employees =?, registered =? WHERE id =?",
		comp.Name,
		comp.Description,
		comp.Type,
		comp.AmountOfEmployees,
		comp.Registered,
		comp.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *CompanyRepository) GetCompany(ctx context.Context, compID string) (*Company, error) {
	row := c.DB.QueryRow(
		"SELECT id, name, description, type, amount_of_employees, registered FROM companies WHERE id =?",
		compID,
	)
	var company Company
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Type, &company.AmountOfEmployees, &company.Registered)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (c *CompanyRepository) Delete(ctx context.Context, compID string) error {
	res, err := c.DB.Exec(
		"DELETE FROM companies WHERE id =?",
		compID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

var CompanyTypes = [4]CompanyType{
	Cooperative,
	Corporations,
	NonProfit,
	SoleProprietorship,
}
