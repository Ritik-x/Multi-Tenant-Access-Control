package repository

import (
	"context"
	"team-access-control/internal/models"
	// "github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepository struct {
	db DBTX
}

func (r *MemberRepository) CreateUser(ctx context.Context, user *models.User) any {
	panic("unimplemented")
}

func NewMembershipRepository(db DBTX) *MemberRepository {
	return &MemberRepository{
		db: db,
	}
}

func (r *MemberRepository) CreateMembership(
	ctx context.Context,
	userID string,
	organizationID string,
	roleID string,
) error {

	query := `
		INSERT INTO memberships (
			user_id,
			organization_id,
			role_id
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		organizationID,
		roleID,
	)
	return err

}

func (r *MemberRepository) GetOrganizationByUserId(ctx context.Context, userID string) (string, error) {
	query := ` SELECT  organization_id  FROM memberships WHERE user_id = $1
		ORDER BY created_at
		LIMIT 1 `
	var organizationID string

	err := r.db.QueryRow(
		ctx, query, userID,
	).Scan(&organizationID)
	if err != nil {
		return "", err
	}

	return organizationID, nil
}
