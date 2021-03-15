package ucCmnpy

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
)

type CompanyUCImpl struct {
	rep company.CompanyRepI
}

func (uc *CompanyUCImpl) Create(company models.Company) (models.Company, error) {
	return uc.rep.CreateCompany(company)
}

func (uc *CompanyUCImpl) GetFavoriteList(userID string, company models.Company) ([]models.Company, error) {
	return uc.rep.GetFavoriteList(userID, company)
}

func (uc *CompanyUCImpl) AddFavorite(userID string, company models.Company) error {
	return uc.rep.AddFavorite(userID, company)
}

func (uc *CompanyUCImpl) DelFavorite(userID string, company models.Company) error {
	return uc.rep.DelFavorite(userID, company)
}
