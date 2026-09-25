package repository

import "context"

type AuditLogRepo struct {
	db DBTX
}


func NewAuditLogRepository ( db DBTX) *AuditLogRepo{
	return &AuditLogRepo{
		db:db,
	}
}

func(r *AuditLogRepo) CreateAuditLogs(
	ctx context.Context,
	organizationID string,
	userID *string,
	action string,
	resource string,
	resourceID *string,
	metadata []byte,
	ipAddress *string,

)error {
	query := `INSERT INTO audit_logs  (organization_id , user_id , action , resource , resource_id , metadata,
			ip_address) VALUES ($1, $2, $3, $4, $5, $6, $7)`

			_, err := r.db.Exec(
		ctx,
		query,
		organizationID,
		userID,
		action,
		resource,
		resourceID,
		metadata,
		ipAddress,
	)

	return err
}