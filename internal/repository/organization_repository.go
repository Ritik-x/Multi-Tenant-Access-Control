package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
db *pgxpool.Pool
}

func NewOrganizationRepository( db *pgxpool.Pool) *OrganizationRepository { 
	return &OrganizationRepository{
		db:db,
	}
}
func ( r *OrganizationRepository ) CreateOrganizationId(ctx context.Context , name string , slug string ) (string, error) {
	query := ` INSERT INTO organizations (
	name , slug) VALUES ($1 , $2) RETURNING id` 
	var OrganizationId string
	err := r.db.QueryRow(ctx , query , name , slug).Scan(&OrganizationId) 
		if err != nil {
		return "", err
	}
		if err != nil {
		return "", err
	}
	return OrganizationId, nil
}