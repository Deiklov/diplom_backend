package ucCmnpy

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	"github.com/pkg/errors"
)

type CompanyUCImpl struct {
	rep company.CompanyRepI
}

func CreateUseCase(cmpnyRepo_ company.CompanyRepI) company.CompanyUCI {
	return &CompanyUCImpl{
		rep: cmpnyRepo_,
	}
}

func (uc *CompanyUCImpl) Create(company models.Company) (*models.Company, error) {
	return uc.rep.CreateCompany(company)
}

func (uc *CompanyUCImpl) GetFavoriteList(userID string, company models.Company) ([]models.Company, error) {
	return uc.rep.GetFavoriteList(userID, company)
}

func (uc *CompanyUCImpl) AddFavorite(userID string, company models.Company) error {
	if company.ID != "" {
		return uc.rep.AddFavorite(userID, company.ID)
	} else if company.Ticker != "" {
		cmpny, err := uc.rep.GetCompany(company.Ticker)
		if err != nil {
			return errors.Wrap(err, "can't get favorite by ticker")
		}
		return uc.rep.AddFavorite(userID, cmpny.ID)

	} else {
		return errors.New("company doesn't contain ID or Name")
	}
}

func (uc *CompanyUCImpl) DelFavorite(userID string, company models.Company) error {
	if company.ID != "" {
		return uc.rep.DelFavorite(userID, company.ID)
	} else if company.Ticker != "" {
		cmpny, err := uc.rep.GetCompany(company.Ticker)
		if err != nil {
			return errors.Wrap(err, "can't get favorite by ticker")
		}
		return uc.rep.DelFavorite(userID, cmpny.ID)

	} else {
		return errors.New("company doesn't contain ID or Name")
	}
}

func (uc *CompanyUCImpl) GetCompany(slug string) (models.Company, error) {
	return uc.rep.GetCompany(slug)
}
func (uc *CompanyUCImpl) GetAllCompanies() ([]models.Company, error) {
	return uc.rep.GetAllCompanies()
}
