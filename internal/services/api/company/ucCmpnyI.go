package company

import "github.com/Deiklov/diplom_backend/internal/models"

type CompanyUCI interface {
	Create(company models.Company) (*models.Company, error)
	GetFavoriteList(userID string, company models.Company) ([]models.Company, error)
	AddFavorite(userID string, company models.Company) error
	DelFavorite(userID string, company models.Company) error
	SearchCompany(slug string) (models.Company, error)
	//maybe should add sort and pagination
	GetAllCompanies() ([]models.Company, error)
}
