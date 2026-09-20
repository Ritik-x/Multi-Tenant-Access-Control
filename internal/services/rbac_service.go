package services

import "context"

type RBACRepository interface {
	GetUserPermissions(
		ctx context.Context,
		userID string,
		organizationID string,
	) ([]string, error)
}

type RBACService struct {
	repo RBACRepository
}

func NewRBACService(repo RBACRepository) *RBACService {
	return &RBACService{
		repo: repo,
	}
}

func (s *RBACService) HasPermission(
	ctx context.Context,
	userID string,
	organizationID string,
	requiredPermission string,
) (bool, error) {

	permissions, err := s.repo.GetUserPermissions(
		ctx,
		userID,
		organizationID,
	)

	if err != nil {
		return false, err
	}

	for _ , permission := range permissions {
		if permission == requiredPermission {
			return true, nil
		}
	}

	return false, nil
}