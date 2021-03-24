package company

import "github.com/Deiklov/diplom_backend/internal/models"

type CompanyRepI interface {
	CreateCompany(company models.Company) (*models.Company, error)
	GetFavoriteList(userID string, company models.Company) ([]models.Company, error)
	AddFavorite(userID string, company models.Company) error
	DelFavorite(userID string, companyID string) error
	GetCompany(slug string) (models.Company, error)
	GetAllCompanies() ([]models.Company, error)
}


