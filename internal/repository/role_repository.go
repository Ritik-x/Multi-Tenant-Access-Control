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



func ( r *RoleRepository) AssignPermission( ctx  context.Context  , roleID string  , permission [] string) error{
	query := `INSERT INTO role_permissions (role_id, permission_id)  SELECT $1, id
		FROM  permissions
		WHERE name = ANY($2)`

		_ , err:= r.db.Exec(ctx , query , roleID, permission , )
		return err
}