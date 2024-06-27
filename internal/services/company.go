package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/knbr13/company-service-go/internal/repositories"
	"github.com/knbr13/company-service-go/internal/validator"
)

type CompanyService struct {
	Repos *repositories.Repositories
}

func (s *CompanyService) Create(ctx context.Context, comp *repositories.Company) (string, error) {
	if err := validateCompany(comp); err != nil {
		return "", err
	}
	comp.ID = uuid.New().String()

	return comp.ID, s.Repos.Company.Insert(ctx, comp)
}

func (s *CompanyService) Get(ctx context.Context, compId string) (*repositories.Company, error) {
	return s.Repos.Company.GetCompany(ctx, compId)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, comp *repositories.Company) error {
	if err := validateCompany(comp); err != nil {
		return err
	}
	return s.Repos.Company.Update(ctx, comp)
}

func (s *CompanyService) Delete(ctx context.Context, compId string) error {
	return s.Repos.Company.Delete(ctx, compId)
}

func validateCompany(comp *repositories.Company) error {
	v := validator.New()
	v.Check(len(comp.Name) >= 3 && len(comp.Name) <= 15, "name", "must be between 3 and 15 characters")
	v.Check(len(comp.Description) <= 3000, "description", "must not exceed 3000 characters")

	validCompanyTypes := make([]string, 0, len(repositories.CompanyTypes))
	for _, t := range repositories.CompanyTypes {
		validCompanyTypes = append(validCompanyTypes, string(t))
	}
	v.Check(slices.Contains(validCompanyTypes, comp.Type), "type", fmt.Sprintf("must be one of: [%s]", strings.Join(validCompanyTypes, ", ")))
	v.Check(comp.AmountOfEmployees != nil, "amount_of_employees", "required field")
	v.Check(comp.Registered != nil, "registered", "required field")

	if !v.Valid() {
		for k, v := range v.Errors {
			return validator.ValidationError(fmt.Sprintf("%s: %s", k, v))
		}
	}

	return nil
}
