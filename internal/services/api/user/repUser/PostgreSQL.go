package repUser

import (
	"database/sql"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	errOwn "github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	"time"
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
	_, err := urep.goquDb.Insert("users").Cols("id", "email", "name", "password").
		Vals(goqu.Vals{user.ID, user.Email, user.Name, user.Password}).Executor().Exec()
	if err != nil {
		logger.Error(err)
		return errOwn.ErrDbBadOperation
	}
	return nil
}

func (urep *UserStore) GetByID(id string) (*models.User, error) {
	usr := models.User{}
	sql, _, err := urep.goquDb.From("users").
		Select("id", "email", "created_at", "updated_at", "name").
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
func (urep *UserStore) Update(user *models.User) (*models.User, error) {
	usrFromDB := models.User{}
	updateFields := make(map[string]interface{})
	updateFields["updated_at"] = time.Now()

	stmt := urep.goquDb.Update("users").
		Where(goqu.C("id").Eq(user.ID)).
		Returning("id", "email", "created_at", "updated_at", "name")
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	stmt=stmt.Set(updateFields)
	sql, _, err := stmt.ToSQL()
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	err = urep.dbsqlx.QueryRowx(sql).StructScan(&usrFromDB)
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrDbBadOperation
	}
	return &usrFromDB, nil
}

func (urep *UserStore) GetByEmail(email string) (*models.User, error) {
	usr := models.User{}
	sql, _, err := urep.goquDb.From("users").
		Select("id", "email", "created_at", "updated_at", "name").
		Where(goqu.C("email").Eq(email)).ToSQL()
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
