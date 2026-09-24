package services

import (
	"context"
	"fmt"
	"strings"
	"team-access-control/internal/repository"
	"time"
)
type InviatationService struct {
	invitationRepo *repository.InviationRepo
		roleRepo       *repository.RoleRepository
	authService *AuthService


}

func NewInvitationService(invitationRepo *repository.InviationRepo ,	roleRepo       *repository.RoleRepository, authService *AuthService,) *InviatationService{
	return &InviatationService{
		invitationRepo: invitationRepo,
		roleRepo:       roleRepo,

		authService:    authService,
	}
}

func ( s *InviatationService) CreateInvitation(	ctx context.Context,
	organizationID string,
	email string,
	roleID string,) (string , string , error){
		email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", "", fmt.Errorf("email is required")
	}



	if roleID == "" {
		return "", "", fmt.Errorf("role id is required")
	}

	roleBelongsToOrg , err := s.roleRepo.RoleBelongsToOrganization(ctx , roleID , organizationID)
	if err != nil {
	return "", "", fmt.Errorf(
		"validate role: %w",
		err,
	)
}
if !roleBelongsToOrg{
	return "", "", fmt.Errorf(
		"role does not belong to organization",
	)
}

	invitationToken , tokenHash , err := s.authService.GenerateInvitationToken()

	if err != nil {
		return "", "", fmt.Errorf(
			"generate invitation token: %w",
			err,
		)
	}
	expiresAt := time.Now().Add(48 * time.Hour)
_, err = s.invitationRepo.CreateInviatation(
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
		return invitationToken, email, nil
	}
