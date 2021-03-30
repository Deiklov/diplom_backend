package repCmpny

import (
	"database/sql"
	"encoding/json"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	errOwn "github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
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
	attributes, err := json.Marshal(cmpny.Attributes)
	if err != nil {
		return nil, err
	}
	ok, err := rep.goquDb.Insert("companies").
		Cols("id", "name", "description", "country", "attributes", "ipo", "ticker", "logo", "weburl").
		Vals(goqu.Vals{uuid.New().String(), strings.ToUpper(cmpny.Name),
			cmpny.Description, cmpny.Country, attributes, cmpny.IPO,
			cmpny.Ticker, cmpny.Logo, cmpny.Weburl}).
		Returning("id", "ticker", "ipo", "name", "country").Executor().ScanStruct(&cmpnyFromDb)
	if err != nil || !ok {
		logger.Error(err)
		return nil, err
	}
	return &cmpnyFromDb, nil
}

func (rep *CompanyRepImpl) GetFavoriteList(userID string, company models.Company) ([]models.Company, error) {
	sql, _, err := rep.goquDb.From("company_by_users").Join(goqu.T("companies"),
		goqu.On(goqu.Ex{"company_by_users.company_id": goqu.I("companies.id")})).
		Where(goqu.C("user_id").Eq(userID)).ToSQL()
	companies := []models.Company{}
	err = rep.dbsqlx.QueryRowx(sql).StructScan(&companies)
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	return companies, nil
}

func (rep *CompanyRepImpl) AddFavorite(userID string, companyID string) error {
	_, err := rep.goquDb.Insert("company_by_users").Cols("id", "company_id", "user_id").
		Vals(goqu.Vals{uuid.New().String(), companyID, userID}).
		Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}

func (rep *CompanyRepImpl) DelFavorite(userID string, companyID string) error {
	_, err := rep.goquDb.Delete("company_by_users").Where(goqu.Ex{
		"user_id":    userID,
		"company_id": companyID,
	}).Executor().Exec()
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (rep *CompanyRepImpl) GetCompany(slug string) (models.Company, error) {
	sql, _, err := rep.goquDb.From("companies").Select("id", "name", "ipo", "country", "ticker").
		Where(goqu.C("ticker").Eq(strings.ToUpper(slug))).ToSQL()
	companies := models.Company{}
	err = rep.dbsqlx.QueryRowx(sql).StructScan(&companies)
	if err != nil {
		logger.Error(err)
		return models.Company{}, err
	}
	return companies, nil
}

func (rep *CompanyRepImpl) GetAllCompanies() ([]models.Company, error) {
	sql, _, err := rep.goquDb.From("companies").Select("id", "name", "ipo", "country", "ticker").ToSQL()
	companies := make([]models.Company, 0)
	rows, err := rep.dbsqlx.Queryx(sql)
	for rows.Next() {
		var cmpn models.Company

		err = rows.StructScan(&cmpn)
		if err != nil {
			return companies, err
		}
		companies = append(companies, cmpn)
	}
	if err != nil {
		logger.Error(err)
		return []models.Company{}, err
	}
	return companies, nil
}
