package services

import (
	"context"
	"encoding/json"
	"fmt"
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





func rbacPermissionCacheKey (organizationId , userId string) string{
	return fmt.Sprintf(
		"%s:%s:%s",
		rbacPermissionCachePrefix,
		organizationId,
		userId,
	)
}


func ( c *RBACCache) setPermissions( 	ctx context.Context,
	organizationID string,
	userID string,
	permissions []string, ) error{

		key := rbacPermissionCacheKey(organizationID, userID)

		data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	return c.client.Set(
		ctx,
		key,
		data,
		c.ttl,
	).Err()
	}



func ( c *RBACCache) GetPermissions(


	ctx context.Context, organizationId string , userId string ,
) ( [ ] string ,bool , error){

	key := rbacPermissionCacheKey(organizationId , userId)

	data , err := c.client.Get(ctx , key ).Result()


	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}


	var permissions []string
		if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, false, err
	}
	return permissions, true, nil
}