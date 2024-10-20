package application

import (
	"company-ms/internal"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func ValidateCompany(company *internal.Company) error {
	validate := validator.New()
	return validate.Struct(company)
}

type CompanyService struct {
	repo      CompanyRepository
	logger    *zap.Logger
	kafkaProd MessageProducer
}

func NewCompanyService(repo CompanyRepository, logger *zap.Logger, kafkaProd MessageProducer) *CompanyService {
	return &CompanyService{repo, logger, kafkaProd}
}

func (s *CompanyService) Create(company *internal.Company) error {
	err := ValidateCompany(company)
	if err != nil {
		s.logger.Error("Validation failed for company")
		return NewBadParamError(err.Error())
	}
	err = s.repo.GetByName(company.Name)
	if err != nil {
		return NewBadParamError(err.Error())
	}
	err = s.repo.Create(company)
	if err != nil {
		s.logger.Error("Error creating company")
		return NewServerError(err.Error())
	}
	s.logger.Info("Company created successfully", zap.Any("company", company))

	message := []byte("Company created: " + company.Name)
	if err := s.kafkaProd.Produce("company_events", message); err != nil {
		s.logger.Error("Failed to produce Kafka message", zap.Error(err))
	}
	return nil
}

func (s *CompanyService) GetAll() ([]*internal.Company, error) {
	companies, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Error fetching all companies", zap.Error(err))
		return nil, err
	}
	return companies, nil
}

func (s *CompanyService) GetByID(id string) (*internal.Company, error) {
	company, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error fetching company by ID", zap.String("id", id), zap.Error(err))
		return nil, NewNotFoundError(err.Error())
	}
	return company, nil
}

func (s *CompanyService) Update(company *internal.Company) error {
	err := ValidateCompany(company)
	if err != nil {
		s.logger.Error("Validation failed for company", zap.Any("company", company), zap.Error(err))
		return NewBadParamError(err.Error())
	}
	err = s.repo.GetByNameAndId(company.ID, company.Name)
	if err != nil {
		return NewBadParamError(err.Error())
	}
	err = s.repo.Update(company)
	if err != nil {
		s.logger.Error("Error updating company", zap.Any("company", company), zap.Error(err))
		return NewServerError(err.Error())
	}
	s.logger.Info("Company updated successfully", zap.Any("company", company))

	message := []byte("Company updated: " + company.Name)
	if err := s.kafkaProd.Produce("company_events", message); err != nil {
		s.logger.Error("Failed to produce Kafka message", zap.Error(err))
	}

	return nil
}

func (s *CompanyService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Error deleting company", zap.String("id", id), zap.Error(err))
		return NewNotFoundError(err.Error())
	}
	s.logger.Info("Company deleted successfully", zap.String("id", id))

	message := []byte("Company deleted with ID: " + id) // Customize the message as needed
	if err := s.kafkaProd.Produce("company_events", message); err != nil {
		s.logger.Error("Failed to produce Kafka message", zap.Error(err))
	}

	return nil
}
