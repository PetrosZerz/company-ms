package application

import (
	"company-ms/internal"
)

// CompanyRepository interface for company data operations
type CompanyRepository interface {
	Create(company *internal.Company) error
	GetByID(id string) (*internal.Company, error)
	GetAll() ([]*internal.Company, error)
	Update(company *internal.Company) error
	Delete(id string) error
	GetByName(name string) error
	GetByNameAndId(id string, name string) error
}
