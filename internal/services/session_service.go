package services

import (
	"context"
	"fmt"
	"team-access-control/internal/models"
	"team-access-control/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService struct {
	db          *pgxpool.Pool
	sessionRepo *repository.SessionRepository
	authService *AuthService
}

func NewSessionService(db *pgxpool.Pool, sessionRepo *repository.SessionRepository,
	authService *AuthService) *SessionService {
	return &SessionService{
		db:          db,
		sessionRepo: sessionRepo,
		authService: authService,
	}

}

func (s *SessionService) RotateRefreshToken(
	ctx context.Context,
	refreshToken string,
) (string, string, string, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"begin transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	txSessionRepo := repository.NewSessionRepository(tx)

	// 1. Hash the refresh token received from client
	refreshTokenHash := s.authService.HashRefreshToken(
		refreshToken,
	)

	// 2. Find and lock the active session
	sessionID, userID, organizationID, expiresAt, err :=
		txSessionRepo.GetSessionByRefreshTokenHash(
			ctx,
			refreshTokenHash,
		)

	if err != nil {
		return "", "", "", fmt.Errorf(
			"get session: %w",
			err,
		)
	}

	// 3. Check expiry
	if time.Now().After(expiresAt) {
		return "", "", "", fmt.Errorf(
			"refresh token expired",
		)
	}

	// 4. Revoke old session
	if err := txSessionRepo.RevokeSession(
		ctx,
		sessionID,
	); err != nil {
		return "", "", "", fmt.Errorf(
			"revoke old session: %w",
			err,
		)
	}

	// 5. Generate new refresh token
	newRefreshToken, newRefreshTokenHash, err :=
		s.authService.GenerateRefreshToken()

	if err != nil {
		return "", "", "", fmt.Errorf(
			"generate refresh token: %w",
			err,
		)
	}

	// 6. Give new refresh token a new expiry
	newExpiresAt := time.Now().Add(
		30 * 24 * time.Hour,
	)

	// 7. Create new session
	_, err = txSessionRepo.CreateSession(
		ctx,
		userID,
		organizationID,
		newRefreshTokenHash,
		newExpiresAt,
	)

	if err != nil {
		return "", "", "", fmt.Errorf(
			"create new session: %w",
			err,
		)
	}

	// 8. Commit everything
	if err := tx.Commit(ctx); err != nil {
		return "", "", "", fmt.Errorf(
			"commit transaction: %w",
			err,
		)
	}

	return newRefreshToken, userID, organizationID, nil
}

func (s *SessionService) GetActiveSession(ctx context.Context,
	userID string) ([]models.Session, error) {
	return s.sessionRepo.GetActiveSessionsByUserID(ctx,
		userID,
	)
}


func ( s *SessionService) RevokeSession( ctx context.Context , sessionID string,
	userID string,) error {
		return s.sessionRepo.RevokeSessionByUserId(
		ctx,
		sessionID,
		userID,
	)
	}