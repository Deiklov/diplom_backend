package repCmpny

import (
	"context"
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/internal/common"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/bxcodec/faker/v3"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"testing"
)

type TestSuite struct {
	suite.Suite
	rep           company.CompanyRepI
	finnhubClient *finnhub.DefaultApiService
	FhubCtx       context.Context
	CmpnyMapper   common.CmpnyHelper
}

func (s *TestSuite) SetupTest() {
	connectionString := fmt.Sprintf("postgres://andrey:167839@localhost:5432/back_db?sslmode=disable", )
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	s.rep = CreateRepCmpny(pdb)
	s.finnhubClient = finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	s.FhubCtx = context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "c0ilbh748v6ot9ddgc0g",
	})

}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCmpnyCreate() {
	randomTicker := []string{"MSFT", "NVDA"}
	profileCmpny, _, err := s.finnhubClient.CompanyProfile2(s.FhubCtx,
		&finnhub.CompanyProfile2Opts{Symbol: optional.NewString(randomTicker[rand.Intn(len(randomTicker))])})
	cmpnModel := s.CmpnyMapper.FinhubProfileToModel(profileCmpny)
	cmpnModel.ID = faker.UUIDHyphenated()
	cmpany, err := s.rep.CreateCompany(cmpnModel)
	s.NotNil(cmpany)
	s.NotNil(cmpany.ID)
	s.Nil(err)
}

func (s *TestSuite) TestAddFavorite() {
	err := s.rep.AddFavorite("1b9405dd-5234-4434-ab93-0ae54d31f25d", models.Company{
		ID: "dafd9393-5552-4b5e-9061-35868031110b",
	})
	s.Nil(err)
}
func (s *TestSuite) TestGetCompany() {
	cmpny, err := s.rep.GetCompany("ABNB")
	s.NotNil(cmpny)
	s.NotEmpty(cmpny.ID)
	s.NotEmpty(cmpny.Ticker)
	s.NotEmpty(cmpny.Ticker)
	s.NotEmpty(cmpny.IPO)
	s.Nil(err)
}

func (s *TestSuite) TestDelFavorite() {
	err := s.rep.DelFavorite("1b9405dd-5234-4434-ab93-0ae54d31f25d", "dafd9393-5552-4b5e-9061-35868031110b")
	s.Nil(err)
}

func (s *TestSuite) TestGetFavorites() {
	companyList, err := s.rep.GetFavoriteList("1b9405dd-5234-4434-ab93-0ae54d31f25d", models.Company{})
	s.Nil(err)
	s.NotNil(companyList)
}
