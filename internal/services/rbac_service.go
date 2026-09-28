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
	cache *RBACCache
}

func NewRBACService(repo RBACRepository , cache *RBACCache,) *RBACService {
	return &RBACService{
		repo: repo,
			cache: cache,
	}
}
func (s *RBACService) HasPermission(
	ctx context.Context,
	userID string,
	organizationID string,
	requiredPermission string,
) (bool, error) {

	permissions, found, err := s.cache.GetPermissions(
		ctx,
		organizationID,
		userID,
	)

	if err != nil {
		// Redis error → PostgreSQL fallback
		permissions, err = s.repo.GetUserPermissions(
			ctx,
			userID,
			organizationID,
		)

		if err != nil {
			return false, err
		}

		// DB permissions → Redis
		_ = s.cache.setPermissions(
			ctx,
			organizationID,
			userID,
			permissions,
		)

	} else if !found {
		// Cache MISS → PostgreSQL
		permissions, err = s.repo.GetUserPermissions(
			ctx,
			userID,
			organizationID,
		)

		if err != nil {
			return false, err
		}

		// DB permissions → Redis
		_ = s.cache.setPermissions(
			ctx,
			organizationID,
			userID,
			permissions,
		)
	}

	// Cache HIT ya DB se permissions milne ke baad
	// dono cases mein yahi permission check chalega.
	for _, permission := range permissions {
		if permission == requiredPermission {
			return true, nil
		}
	}

	return false, nil
}