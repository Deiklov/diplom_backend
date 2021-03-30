package common

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/bxcodec/faker/v3"
	"time"
)

type CmpnyHelper struct {
}

func (cpny *CmpnyHelper) FinhubProfileToModel(profile2 finnhub.CompanyProfile2) models.Company {
	t, err := time.Parse("2006-01-02", profile2.Ipo)
	if err != nil {
		return models.Company{}
	}

	return models.Company{
		ID:   "",
		Name: profile2.Name,
		IPO:  t,
		//нету нормально описания
		Description: faker.Sentence(),
		Country:     profile2.Country,
		Ticker:      profile2.Ticker,
		Logo:        profile2.Logo,
		Weburl:      profile2.Weburl,
		Attributes: models.AttributesCmpny{
			Currency: profile2.Currency,
			Exchange: profile2.Exchange,
			Industry: profile2.FinnhubIndustry,
		},
	}
}
