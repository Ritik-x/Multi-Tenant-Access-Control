package repository

import (
	"context"
	"time"
)

type InviationRepo struct {
	db DBTX
}

func NewInviationRepository(db DBTX) *InviationRepo {
	return &InviationRepo{
		db: db,
	}
}

func (r *InviationRepo) CreateInviatation(
	ctx context.Context,
	organizationID string,
	email string,
	roleID string,
	tokenHash string,
	expiresAt time.Time,
) (string, error) {
	query := `INSERT INTO (organization_id ,email , role_id , token_hash , expires_at ) VALUES ($1 , $2 , $3, $4, $5 ) RETURNING id`

	var invitationID string
	err := r.db.QueryRow(ctx, query, organizationID,
		email,
		roleID,
		tokenHash,
		expiresAt).Scan(&invitationID)

	if err != nil {
		return "", err
	}
	return invitationID, nil

}
