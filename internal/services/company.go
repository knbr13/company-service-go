package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/knbr13/company-service-go/internal/repositories"
	"github.com/knbr13/company-service-go/internal/validator"
)

const (
	CompanyCreated string = "company_created"
	CompanyUpdated string = "company_updated"
	CompanyDeleted string = "company_deleted"
)

type CompanyService struct {
	Repos    *repositories.Repositories
	Producer sarama.SyncProducer
	ErrCh    chan<- error
}

func (s *CompanyService) Create(ctx context.Context, comp *repositories.Company) (string, error) {
	if err := validateCompany(comp); err != nil {
		return "", err
	}
	comp.ID = uuid.New().String()

	if err := s.Repos.Company.Insert(ctx, comp); err != nil {
		return "", err
	}

	go func(producer sarama.SyncProducer, comp *repositories.Company, errCh chan<- error) {
		event := map[string]interface{}{
			"event":   CompanyCreated,
			"company": comp,
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal event: %w", err)
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: CompanyCreated,
			Key:   sarama.StringEncoder(comp.ID),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = producer.SendMessage(msg)
		if err != nil {
			errCh <- fmt.Errorf("failed to send message: %w", err)
			return
		}

		log.Printf("Message sent to topic %s for company %s", CompanyCreated, comp.ID)
	}(s.Producer, comp, s.ErrCh)

	return comp.ID, nil
}

func (s *CompanyService) Get(ctx context.Context, compId string) (*repositories.Company, error) {
	return s.Repos.Company.GetCompany(ctx, compId)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, comp *repositories.Company) error {
	if err := validateCompany(comp); err != nil {
		return err
	}

	if err := s.Repos.Company.Update(ctx, comp); err != nil {
		return err
	}

	go func(producer sarama.SyncProducer, comp *repositories.Company, errCh chan<- error) {
		event := map[string]interface{}{
			"event":   CompanyUpdated,
			"company": comp,
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal event: %w", err)
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: CompanyUpdated,
			Key:   sarama.StringEncoder(comp.ID),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = producer.SendMessage(msg)
		if err != nil {
			errCh <- fmt.Errorf("failed to send message: %w", err)
			return
		}

		log.Printf("Message sent to topic %s for company %s", CompanyUpdated, comp.ID)
	}(s.Producer, comp, s.ErrCh)

	return nil
}

func (s *CompanyService) Delete(ctx context.Context, compId string) error {
	if err := s.Repos.Company.Delete(ctx, compId); err != nil {
		return err
	}

	go func(producer sarama.SyncProducer, compId string, errCh chan<- error) {
		event := map[string]interface{}{
			"event":      CompanyDeleted,
			"company_id": compId,
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal event: %w", err)
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: CompanyDeleted,
			Key:   sarama.StringEncoder(compId),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = producer.SendMessage(msg)
		if err != nil {
			errCh <- fmt.Errorf("failed to send message: %w", err)
			return
		}

		log.Printf("Message sent to topic %s for company %s", CompanyDeleted, compId)
	}(s.Producer, compId, s.ErrCh)

	return nil
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
		var sb strings.Builder
		for k, v := range v.Errors {
			sb.WriteString(fmt.Sprintf("%s: %s; ", k, v))
		}
		return validator.ValidationError(sb.String())
	}

	return nil
}
