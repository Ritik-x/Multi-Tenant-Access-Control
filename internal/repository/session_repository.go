package repository

import (
	"context"
	"team-access-control/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type SessionRepository struct {
	db DBTX
}

func NewSessionRepository(db DBTX) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
	userId string,
	organizationId string,
	refreshTokenHash string,
	expiresAt time.Time,
) (string, error) {
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

	return sessionID, err
}

func (r *SessionRepository) GetSessionByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (string, string, string, time.Time, error) {
	query := `
	SELECT
		id,
		user_id,
		organization_id,
		expires_at
	FROM sessions
	WHERE refresh_token_hash = $1
	  AND revoked_at IS NULL
	  		FOR UPDATE
`

	var (
		sessionID      string
		userId         string
		organizationId string
		expiresAt      time.Time
	)
	err := r.db.QueryRow(ctx,
		query,
		refreshTokenHash).Scan(&sessionID,
		&userId,
		&organizationId,
		&expiresAt)
	if err != nil {
		return "", "", "", time.Time{}, err
	}
	return sessionID, userId, organizationId, expiresAt, nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, sessionId string) error {
	query := ` UPDATE sessions 	SET revoked_at = NOW() WHERE id = $1
		  AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx,
		query,
		sessionId)

	return err
}

func (r *SessionRepository) GetActiveSessionsByUserID(
	ctx context.Context,
	userID string,
) ([]models.Session, error) {

	query := `
		SELECT
			id,
			user_id,
			organization_id,
			expires_at,
			created_at,
			revoked_at
		FROM sessions
		WHERE user_id = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session

	for rows.Next() {
		var session models.Session

		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.OrganizationID,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.RevokedAt,
		); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) RevokeSessionByUserId(ctx context.Context, sessionId string, userId string) error {
	query := `UPDATE sessions SET revoked_at = NOW()
		WHERE id = $1
		  AND user_id = $2
		  AND revoked_at IS NULL`

	tag, err := r.db.Exec(ctx,
		query,
		sessionId,
		userId,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *SessionRepository) RevokeAllSessions(ctx context.Context, userID string) error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1
		  AND revoked_at IS NULL
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
	)

	return err
}
