package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/user"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	// create
	Create(c context.Context, tx pgx.Tx, user *user.User) error

	// read
	FindAll(c context.Context, filter *model.QueryParam) ([]user.User, error)
	Count(c context.Context, filter *model.QueryParam) (int, error)
	FindByColumn(c context.Context, column string, value interface{}) (*user.User, error)

	// update
	Update(c context.Context, tx pgx.Tx, user *user.User) error

	// delete
	SoftDelete(c context.Context, tx pgx.Tx, NIKs ...string) error
	Delete(c context.Context, tx pgx.Tx, NIKs ...string) error
}

type userRepository struct {
	Logger *zap.SugaredLogger
	DB     *pgxpool.Pool
}

func New(
	logger *zap.SugaredLogger,
	db *pgxpool.Pool,
) UserRepository {
	return &userRepository{
		Logger: logger,
		DB:     db,
	}
}

func (ur *userRepository) Create(c context.Context, tx pgx.Tx, user *user.User) error {
	queryStr := `
	INSERT INTO users (
		id,
		full_name,
		email,
		password,
		role,
		last_activity_date,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err := tx.Exec(c, queryStr,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.LastActivityDate,
		user.CreatedAt,
	); err != nil {
		ur.Logger.Errorw("failed to create user", "error", err)
		return err
	}

	return nil
}

func (ur *userRepository) FindAll(c context.Context, filter *model.QueryParam) ([]user.User, error) {
	queryStr := `
	SELECT 
		id,
		full_name,
		email,
		role,
		last_activity_date,
		created_at
	FROM users`

	// filter
	utils.Json(filter)
	queryStr, args := filterUsers(queryStr, filter)

	// sort
	queryStr, err := query.Sort(queryStr, filter.Sort, SortUserMap)
	if err != nil {
		ur.Logger.Errorw("failed to sort query", "error", err)
		return nil, err
	}

	// pagination
	queryStr = query.Paginate(queryStr, filter.Page, filter.Limit)

	query.Debug(queryStr, args...)
	rows, err := ur.DB.Query(c, queryStr, args...)
	if err != nil {
		ur.Logger.Errorw("Error on find users", "error", err.Error())
		return nil, err
	}

	defer rows.Close()

	var users []user.User
	for rows.Next() {
		user := user.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.LastActivityDate,
			&user.CreatedAt,
		); err != nil {
			ur.Logger.Errorw("error on Find User", "error", err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *userRepository) Count(c context.Context, filter *model.QueryParam) (int, error) {
	queryStr := `
	SELECT 
		COUNT(u.id)
	FROM users u`

	queryStr, args := filterUsers(queryStr, filter)

	var count int
	if err := ur.DB.QueryRow(c, queryStr, args...).Scan(&count); err != nil {
		ur.Logger.Errorw("error on count user", "error", err.Error())
		return 0, err
	}

	return count, nil
}

func (ur *userRepository) FindByColumn(c context.Context, column string, value interface{}) (*user.User, error) {
	queryStr := fmt.Sprintf(`
	SELECT 
		id,
		full_name,
		email,
		password,
		role,
		last_activity_date,
		created_at
	FROM users
	WHERE %s = $1`, column)

	var user user.User
	if err := ur.DB.QueryRow(c, queryStr, value).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.LastActivityDate,
		&user.CreatedAt,
	); err != nil {
		utils.Debug(err)
		ur.Logger.Errorw("error on find User by NIK ", "error", err.Error())
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) Update(c context.Context, tx pgx.Tx, user *user.User) error {
	queryStr := `
	UPDATE users SET
		full_name = $2,
		email = $3,
		password = $4,
		role = $5,
		last_activity_date = $6
	WHERE id = $1`

	if _, err := tx.Exec(c, queryStr,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.LastActivityDate,
	); err != nil {
		ur.Logger.Errorw("failed to update user", "error", err)
		return err
	}

	return nil
}

func (ur *userRepository) SoftDelete(c context.Context, tx pgx.Tx, NIKs ...string) error {
	queryStr := `
	UPDATE users SET
		deleted_at = $2
	WHERE nik = ANY($1)`

	if _, err := tx.Exec(c, queryStr,
		NIKs,
		time.Now(),
	); err != nil {
		ur.Logger.Errorw("failed to delete user", "error", err)
		return err
	}

	return nil
}

func (ur *userRepository) Delete(c context.Context, tx pgx.Tx, NIKs ...string) error {
	queryStr := `
	DELETE FROM users
	WHERE nik = ANY($1)`

	if _, err := tx.Exec(c, queryStr, NIKs); err != nil {
		ur.Logger.Errorw("failed to delete user", "error", err)
		return err
	}

	return nil
}
