package repUser

import (
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/bxcodec/faker/v3"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type TestSuite struct {
	suite.Suite
	userrep user.Repository
}

func (s *TestSuite) SetupTest() {
	connectionString := fmt.Sprintf("postgres://andrey:167839@localhost:5432/back_db?sslmode=disable", )
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	s.userrep = CreateRepository(pdb)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (s *TestSuite) TestUserCreate() {
	err := s.userrep.Create(&models.User{
		ID:    faker.UUIDHyphenated(),
		Name:  faker.Name(),
		Email: faker.Email(),
	})
	s.Nil(err)
}
func (s *TestSuite) TestUserGet() {
	userID := faker.UUIDHyphenated()
	userName := faker.Name()
	userEmail := faker.Email()
	err := s.userrep.Create(&models.User{
		ID:    userID,
		Name:  userName,
		Email: userEmail,
	})
	userDB, err := s.userrep.GetByID(userID)
	s.NotNil(userDB)
	s.Equal(userDB.ID, userID)
	s.Equal(userDB.Name, userName)
	s.Equal(userDB.Email, userEmail)

	s.Nil(err)
}
func (s *TestSuite) TestUserUpdate() {
	userID := "044ae442-ae48-4096-b2b8-5809733089e1"
	userName := faker.Name()
	userEmail := faker.Email()
	usrFromDB, err := s.userrep.Update(&models.User{
		ID:    userID,
		Name:  userName,
		Email: userEmail,
	})
	s.NotNil(usrFromDB)
	s.Equal(userEmail, usrFromDB.Email)
	s.Equal(userName, usrFromDB.Name)

	s.Nil(err)
}
