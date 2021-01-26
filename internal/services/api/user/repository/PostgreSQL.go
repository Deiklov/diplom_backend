package repository

import (
	"database/sql"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	errOwn "github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	DB     *sql.DB
	goquDb *goqu.Database
	dbsqlx *sqlx.DB
}

func CreateRepository(db *sql.DB) user.Repository {
	return &UserStore{
		DB:     db,
		dbsqlx: sqlx.NewDb(db, "postgres"),
		goquDb: goqu.New("postgres", db),
	}
}

func (urep *UserStore) Create(user *models.User) error {
	_, err := urep.goquDb.Insert("users").Cols("id", "phone", "name").
		Vals(goqu.Vals{user.ID, user.Phone, user.Name}).Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}

func (urep *UserStore) GetByID(id string) (*models.User, error) {
	usr := models.User{}
	sql, _, err := urep.goquDb.From("users").
		Select("id", "phone", "created_at", "updated_at", "name").
		Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	//maybe reflect
	err = urep.dbsqlx.QueryRowx(sql).StructScan(&usr)
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	return &usr, nil
}

func (urep *UserStore) Delete(id string) error {
	_, err := urep.goquDb.Delete("users").Where(goqu.C("id").Eq(id)).Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}
