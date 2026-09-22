package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"team-access-control/internal/models"
	"team-access-control/internal/repository"
)

type RegistrationService struct {
	db             *pgxpool.Pool
	userRepo       *repository.UserRepository
	orgRepo        *repository.OrganizationRepository
	roleRepo       *repository.RoleRepository
	membershipRepo *repository.MemberRepository
}

func NewRegistrationService(
	db *pgxpool.Pool,
	userRepo *repository.UserRepository,
	orgRepo *repository.OrganizationRepository,
	roleRepo *repository.RoleRepository,
	membershipRepo *repository.MemberRepository,
) *RegistrationService {
	return &RegistrationService{
		db:             db,
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		roleRepo:       roleRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *RegistrationService) Register(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
	orgName string,
	orgSlug string,
) (*models.User, string, string, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txUserRepo := repository.NewUserRepository(tx)
	txOrgRepo := repository.NewOrganizationRepository(tx)
	txRoleRepo := repository.NewRoleRepository(tx)
	txMembershipRepo := repository.NewMembershipRepository(tx)

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := txUserRepo.CreateUser(ctx, user); err != nil {
		return nil, "", "", fmt.Errorf("create user: %w", err)
	}

	organizationID, err := txOrgRepo.CreateOrganizationId(
		ctx,
		orgName,
		orgSlug,
	)
	if err != nil {
		return nil, "", "", fmt.Errorf("create organization: %w", err)
	}

	roleID, err := txRoleRepo.CreateRole(
		ctx,
		organizationID,
		"owner",
	)
	if err != nil {
		return nil, "", "", fmt.Errorf("create owner role: %w", err)
	}

	//  Assign al default permisins t he owner  role
	ownerPermissions := [] string{
		"users.read",
	"users.create",
	"users.update",
	"users.delete",
	"team.read",
	"team.invite",
	"team.remove",
	"audit.read",
	}

	if err :=txRoleRepo.AssignPermission(	ctx,
	roleID,
	ownerPermissions,);err !=nil{
		return nil, "", "", fmt.Errorf("assign owner permissions: %w", err)
	}
	if err := txMembershipRepo.CreateMembership(
		ctx,
		user.ID,
		organizationID,
		roleID,
	); err != nil {
		return nil, "", "", fmt.Errorf("create membership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", "", fmt.Errorf("commit transaction: %w", err)
	}

	return user, organizationID, roleID, nil
}
