package repository

import (
	"context"
	// "github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db DBTX
}

func NewRoleRepository(db DBTX) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) CreateRole(ctx context.Context, organizationID string, name string) (string, error) {
	query := `INSERT INTO roles (organization_id,
			name) VALUES ($1, $2)
		RETURNING id `
	var roleID string
	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		name,
	).Scan(&roleID)
	if err != nil {
		return "", err
	}
	return roleID, nil
}
