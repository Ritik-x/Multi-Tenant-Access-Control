package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Membership struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	OrganizationID string    `json:"organization_id"`
	RoleID         string    `json:"role_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Role struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}

type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
type Session struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	OrganizationID string     `json:"organization_id"`
	ExpiresAt      time.Time  `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
}

type Invitation struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Email          string     `json:"email"`
	RoleID         string     `json:"role_id"`
	ExpiresAt      time.Time  `json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
		CreatedAt      time.Time              `json:"created_at"`
}



type AuditLogs struct{
	ID string `json:"id"`
	OrganizationID string     `json:"organization_id"`
	UserID         *string    `json:"user_id"`
	Action string  `json:"action"`
	Resource       string                 `json:"resource"`
	ResourceID     *string                `json:"resource_id,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	IPAddress      *string                `json:"ip_address,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}