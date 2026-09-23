package repository

import (
	"context"
	"time"
)

type SessionRepository struct {
	db DBTX
}
func NewSessionRepository (db DBTX) *SessionRepository {
	return &SessionRepository{db: db}
}


func ( r *SessionRepository ) CreateSession ( 
	ctx context.Context ,
	userId string ,
	organizationId string ,
	refreshTokenHash string,
	expiresAt time.Time,
) (string , error) {
	query := `INSERT INTO sessions (
	user_id , organization_id,
			refresh_token_hash,
			expires_at)	VALUES ($1, $2, $3, $4)
		RETURNING id`
		var sessionID string
		err := r.db.QueryRow(ctx,
		query,
		userId,
		organizationId,
		refreshTokenHash,
		expiresAt).Scan(&sessionID)

		return sessionID ,err
}

func (r *SessionRepository  ) GetSessionByRefreshTokenHash( ctx context.Context , refreshTokenHash string ) (string, string, string, time.Time, error){
	query:= `SELECT id , user_id , organization_id , expires_at  FROM sessions WHERE efresh_token_hash = $1
		  AND revoked_at IS NULL`

		  var (sessionID string 
		userId string 
		organizationId string 
			expiresAt      time.Time

		)
		err := r.db.QueryRow(ctx,
		query,
		refreshTokenHash,).Scan(&sessionID,
		&userId,
		&organizationId,
		&expiresAt,)
		if err != nil {
		return "", "", "", time.Time{}, err
	}
		return sessionID, userId, organizationId, expiresAt, nil
}

func ( r *SessionRepository) RevokeSession (ctx context.Context , sessionId string) error {
	query := ` UPDATE sessions 	SET revoked_at = NOW() WHERE id = $1
		  AND revoked_at IS NULL`
		  _ , err := r.db.Exec(ctx,
		query,
		sessionId,)
		
		return err
}


