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


func ( r *InviationRepo) GetInvitationByTokenHash(ctx context.Context , tokenHash string)(string , string , string , string , time.Time , *time.Time , error){

	query := `SELECT id , organization_id email,
			role_id,
			expires_at,
			accepted_at
		FROM invitations
		WHERE token_hash = $1`

		var (
			invitationID string
			organizationID string 
			email string 
			roleID         string
		expiresAt      time.Time
		acceptedAt     *time.Time
		)
		err := r.db.QueryRow(ctx , query , tokenHash).Scan(&invitationID,
		&organizationID,
		&email,
		&roleID,
		&expiresAt,
		&acceptedAt,)
		if err != nil {
		return "", "", "", "", time.Time{}, nil, err
	}
		return invitationID, organizationID, email, roleID, expiresAt, acceptedAt, nil

}





func (r *InviationRepo) GetInvitationByTokenHashForUpdate(
	ctx context.Context,
	tokenHash string,
) (
	string,
	string,
	string,
	string,
	time.Time,
	*time.Time,
	error,
) {

	query := `
		SELECT
			id,
			organization_id,
			email,
			role_id,
			expires_at,
			accepted_at
		FROM invitations
		WHERE token_hash = $1
		FOR UPDATE
	`

	var (
		invitationID   string
		organizationID string
		email          string
		roleID         string
		expiresAt      time.Time
		acceptedAt     *time.Time
	)

	err := r.db.QueryRow(
		ctx,
		query,
		tokenHash,
	).Scan(
		&invitationID,
		&organizationID,
		&email,
		&roleID,
		&expiresAt,
		&acceptedAt,
	)

	if err != nil {
		return "", "", "", "", time.Time{}, nil, err
	}

	return invitationID, organizationID, email, roleID, expiresAt, acceptedAt, nil
}

func ( r *InviationRepo) MarkInvitationAccepted(ctx context.Context , tokenHash string) error {


	query := `UPDATE invitations SET accepted_at = NOW() WHERE token_hash = $1
		  AND accepted_at IS NULL`
		  	_, err := r.db.Exec(
		ctx,
		query,
		tokenHash,
	)

	return err
}
