package services

import (
	"context"
	"encoding/json"
	"fmt"
	"team-access-control/internal/models"
	"team-access-control/internal/repository"
)

type AuditLogService struct{
	repo *repository.AuditLogRepo
}

func NewAuditLogService( repo *repository.AuditLogRepo,)*AuditLogService{
	return &AuditLogService{
		repo : repo,

	}
}

func ( s *AuditLogService) Logs(
	ctx context.Context,
	organizationID string,
	userID *string,
	action string,
	resource string,
	resourceID *string,
	metadata map[string]interface{},
	ipAddress *string,
)error{
	var metadataJSON []byte 

	if metadata != nil {
		data , err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("marshal audit metadata: %w", err)
		}

		metadataJSON = data
	}
	return s.repo.CreateAuditLogs(ctx,
		organizationID,
		userID,
		action,
		resource,
		resourceID,
		metadataJSON,
		ipAddress,)
}


func( s *AuditLogService) GetLogs ( ctx context.Context , organizationId string)  ([]models.AuditLogs , error){
	logs , err := s.repo.GetAuditLogs(
		ctx , organizationId,
	)
		if err != nil {
		return nil, fmt.Errorf("get audit logs: %w", err)
	}
	return logs , nil
}




