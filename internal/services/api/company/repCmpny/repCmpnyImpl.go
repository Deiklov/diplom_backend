package repCmpny

import (
	"database/sql"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	errOwn "github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CompanyRepImpl struct {
	DB     *sql.DB
	goquDb *goqu.Database
	dbsqlx *sqlx.DB
}

func CreateRepCmpny(db *sql.DB) company.CompanyRepI {
	return &CompanyRepImpl{
		DB:     db,
		dbsqlx: sqlx.NewDb(db, "postgres"),
		goquDb: goqu.New("postgres", db),
	}
}

func (rep *CompanyRepImpl) CreateCompany(cmpny models.Company) (*models.Company, error) {
	cmpnyFromDb := models.Company{}
	ok, err := rep.goquDb.Insert("companies").Cols("id", "name", "year", "description").
		Vals(goqu.Vals{cmpny.ID, cmpny.Name, cmpny.Year, cmpny.Description}).
		Returning("id", "name", "year", "description").Executor().ScanStruct(&cmpnyFromDb)
	if err != nil || !ok {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	return &cmpnyFromDb, nil
}

func (rep *CompanyRepImpl) GetFavoriteList(userID string, company models.Company) ([]*models.Company, error) {
	sql, _, err := rep.goquDb.From("company_by_users").Join(goqu.T("companies"),
		goqu.On(goqu.Ex{"company_by_users.company_id": goqu.I("companies.id")})).
		Where(goqu.C("user_id").Eq(userID)).ToSQL()
	companies := []models.Company{}
	err = rep.dbsqlx.QueryRowx(sql).StructScan(&companies)
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	return nil, nil
}

func (rep *CompanyRepImpl) AddFavorite(userID string, company models.Company) error {
	_, err := rep.goquDb.Insert("company_by_users").Cols("id", "company_id", "user_id").
		Vals(goqu.Vals{uuid.New().String(), company.ID, userID}).
		Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}

func (rep *CompanyRepImpl) DelFavorite(userID string, company models.Company) error {
	_, err := rep.goquDb.Delete("company_by_users").Where(goqu.Ex{
		"user_id":    userID,
		"company_id": company.ID,
	}).Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}
