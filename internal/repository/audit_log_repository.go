package repository

import (
	"context"
	"encoding/json"
	"team-access-control/internal/models"
)

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

func ( r*AuditLogRepo) GetAuditLogs(ctx context.Context , organizationId string ) ([]models.AuditLogs,error){
	query := `SELECT id , organization_id , user_id , action , resource , resource_id , metadata,
			ip_address,
			created_at FROM audit_logs   WHERE organization_id = $1
		ORDER BY created_at DESC`

		rows , err := r.db.Query(ctx , query , organizationId )
			if err != nil {
		return nil, err
	}
	defer rows.Close()


	var logs []models.AuditLogs
	for rows.Next(){
		var log models.AuditLogs
		var metaData []byte
		var ipAddress *string
			err := rows.Scan(
			&log.ID,
			&log.OrganizationID,
			&log.UserID,
			&log.Action,
			&log.Resource,
			&log.ResourceID,
			&metaData,
			&ipAddress,
			&log.CreatedAt,
		) 
		if err !=nil {
			return nil , err
		}
		log.IPAddress = ipAddress

		if metaData != nil {
			if err := json.Unmarshal(metaData , &log.Metadata); err != nil {
return nil , err
			}
		}
		logs = append(logs, log)
		
	}
		if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs , err
}