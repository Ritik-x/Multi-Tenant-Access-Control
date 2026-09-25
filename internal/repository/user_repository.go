package repository

import (
	"context"
	"team-access-control/internal/models"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func ( r *UserRepository) CreateUser (
	ctx context.Context , 
user *models.User,
) error {

	query := ` INSERT INTO users (
	name,
			email,
			password_hash
			
			) VALUES  ($1, $2, $3)
			 
			RETURNING
			id,
			name,
			email,
			password_hash,
			created_at,
			updated_at
			`

			return r.db.QueryRow(
				ctx , query , user.Name , user.Email 	,user.PasswordHash,	).Scan(&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,)
}
func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}