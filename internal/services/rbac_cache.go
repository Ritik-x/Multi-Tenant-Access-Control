package services

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const rbacPermissionCachePrefix = "rbac:permissions"

type RBACCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRBACCache(
	client *redis.Client,
	ttl time.Duration,
) *RBACCache {
	return &RBACCache{
		client: client,
		ttl:    ttl,
	}
}
