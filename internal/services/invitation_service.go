package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"team-access-control/internal/repository"
)

type InvitationService struct {
	db              *pgxpool.Pool
	invitationRepo  *repository.InviationRepo
	roleRepo        *repository.RoleRepository
	membershipRepo  *repository.MemberRepository
	authService     *AuthService
	auditLogService *AuditLogService
}

func NewInvitationService(
	db *pgxpool.Pool,
	invitationRepo *repository.InviationRepo,
	roleRepo *repository.RoleRepository,
	membershipRepo *repository.MemberRepository,
	authService *AuthService,
	auditLogService *AuditLogService,
) *InvitationService {
	return &InvitationService{
		db:              db,
		invitationRepo:  invitationRepo,
		roleRepo:        roleRepo,
		membershipRepo:  membershipRepo,
		authService:     authService,
		auditLogService: auditLogService,
	}
}

func (s *InvitationService) CreateInvitation(
	ctx context.Context,
	organizationID string,
	userID string,
	email string,
	roleID string,
	ipAddress string,
) (string, string, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return "", "", fmt.Errorf("email is required")
	}

	if roleID == "" {
		return "", "", fmt.Errorf("role id is required")
	}

	roleBelongsToOrg, err := s.roleRepo.RoleBelongsToOrganization(
		ctx,
		roleID,
		organizationID,
	)
	if err != nil {
		return "", "", fmt.Errorf("validate role: %w", err)
	}

	if !roleBelongsToOrg {
		return "", "", fmt.Errorf("role does not belong to organization")
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	// Repositories using the same transaction
	txInvitationRepo := repository.NewInviationRepository(tx)
	// txAuditLogRepo := respository.NewAuditLogRepository(tx)

	invitationToken, tokenHash, err :=
		s.authService.GenerateInvitationToken()

	if err != nil {
		return "", "", fmt.Errorf(
			"generate invitation token: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(48 * time.Hour)

	_, err = txInvitationRepo.CreateInviatation(
		ctx,
		organizationID,
		email,
		roleID,
		tokenHash,
		expiresAt,
	)

	if err != nil {
		return "", "", fmt.Errorf(
			"create invitation: %w",
			err,
		)
	}

	var ip *string

	if ipAddress != "" {
		ip = &ipAddress
	}

	metaData := map[string]interface{}{
		"email":   email,
		"role_id": roleID,
	}

	if err := s.auditLogService.Logs(
		ctx,
		// txAuditLogRepo,
		organizationID,
		&userID,
		"invitation.created",
		"invitation",
		nil,
		metaData,
		ip,
		); err != nil {
		return "", "", fmt.Errorf("create audit log: %w", err)
	}

	// Commit invitation + audit log together
	if err := tx.Commit(ctx); err != nil {
		return "", "", fmt.Errorf("commit transaction: %w", err)
	}

	return invitationToken, email, nil
}

func (s *InvitationService) AcceptInvitation(
	ctx context.Context,
	token string,
	userEmail string,
	userID string,
) error {

	token = strings.TrimSpace(token)
	userEmail = strings.ToLower(strings.TrimSpace(userEmail))

	if token == "" {
		return fmt.Errorf("invitation token is required")
	}

	if userEmail == "" {
		return fmt.Errorf("user email is required")
	}

	if userID == "" {
		return fmt.Errorf("user id is required")
	}

	tokenHash := s.authService.HashRefreshToken(token)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	txInvitationRepo := repository.NewInviationRepository(tx)
	txMembershipRepo := repository.NewMembershipRepository(tx)

	// Lock invitation row
	_, organizationID, invitedEmail, roleID, expiresAt, acceptedAt, err :=
		txInvitationRepo.GetInvitationByTokenHashForUpdate(
			ctx,
			tokenHash,
		)

	if err != nil {
		return fmt.Errorf("get invitation: %w", err)
	}

	if acceptedAt != nil {
		return fmt.Errorf("invitation already accepted")
	}

	if time.Now().After(expiresAt) {
		return fmt.Errorf("invitation expired")
	}

	if userEmail != invitedEmail {
		return fmt.Errorf("invitation email does not match user")
	}

	// Create membership
	if err := txMembershipRepo.CreateMembership(
		ctx,
		userID,
		organizationID,
		roleID,
	); err != nil {
		return fmt.Errorf("create membership: %w", err)
	}

	// Mark invitation as accepted
	if err := txInvitationRepo.MarkInvitationAccepted(
		ctx,
		tokenHash,
	); err != nil {
		return fmt.Errorf("mark invitation accepted: %w", err)
	}

	// Commit membership + invitation update
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}