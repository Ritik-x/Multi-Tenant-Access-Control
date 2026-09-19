package repository

import (
	"context"
	"team-access-control/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository ( db *pgxpool.Pool) *UserRepository{
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