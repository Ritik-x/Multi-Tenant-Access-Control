package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// import "github.com/jackc/pgx/v5"

type RBACrepository struct {
db  *pgxpool.Pool
}

func NewRBACRepository ( db *pgxpool.Pool) * RBACrepository{
	return &RBACrepository{
		db : db,

	}
}

func ( r *RBACrepository) GetUserPermissions (
	ctx context.Context,
	userId string ,
	organizationId string ,


) ([]string, error) {
	query := ` SELECT DISTINCT p.name
		FROM memberships m
		JOIN roles ro
			ON ro.id = m.role_id
		JOIN role_permissions rp
			ON rp.role_id = ro.id
		JOIN permissions p
			ON p.id = rp.permission_id
		WHERE m.user_id = $1
		  AND m.organization_id = $2`
		  rows , err :=r.db.Query(ctx , query , userId , organizationId)
		  if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions [] string
	for rows.Next() {
		var permission string 
		if err := rows.Scan(&permission) ; err != nil {
				return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}