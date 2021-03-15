package repCmpny

import (
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	"github.com/bxcodec/faker/v3"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/suite"
	"log"
	"strconv"
	"testing"
)

type TestSuite struct {
	suite.Suite
	rep company.CompanyRepI
}

func (s *TestSuite) SetupTest() {
	connectionString := fmt.Sprintf("postgres://andrey:167839@localhost:5432/back_db?sslmode=disable", )
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	s.rep = CreateRepCmpny(pdb)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestUserCreate() {
	company, err := s.rep.CreateCompany(models.Company{
		ID:          faker.UUIDHyphenated(),
		Name:        faker.Username(),
		Year:        func() uint32 { data, _ := strconv.Atoi(faker.YearString()); return uint32(data) }(),
		Description: faker.Sentence(),
	})
	s.NotNil(company)
	s.NotNil(company.ID)
	s.Nil(err)
}

func (s *TestSuite) TestAddFavorite() {
	err := s.rep.AddFavorite("1b9405dd-5234-4434-ab93-0ae54d31f25d", models.Company{
		ID: "dafd9393-5552-4b5e-9061-35868031110b",
	})
	s.Nil(err)
}

func (s *TestSuite) TestDelFavorite() {
	err := s.rep.DelFavorite("1b9405dd-5234-4434-ab93-0ae54d31f25d", models.Company{
		ID: "dafd9393-5552-4b5e-9061-35868031110b",
	})
	s.Nil(err)
}

func (s *TestSuite) TestGetFavorites() {
	companyList, err := s.rep.GetFavoriteList("1b9405dd-5234-4434-ab93-0ae54d31f25d", models.Company{})
	s.Nil(err)
	s.NotNil(companyList)
}
