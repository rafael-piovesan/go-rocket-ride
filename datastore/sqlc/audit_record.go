package sqlc

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/audit"
)

func (s *sqlStore) CreateAuditRecord(ctx context.Context, ar *entity.AuditRecord) (*entity.AuditRecord, error) {
	arg := CreateAuditRecordParams{
		Action:       ar.Action.String(),
		CreatedAt:    ar.CreatedAt,
		Data:         ar.Data,
		OriginIp:     ar.OriginIP,
		ResourceID:   ar.ResourceID,
		ResourceType: ar.ResourceType.String(),
		UserID:       ar.UserID,
	}

	model, err := s.q.CreateAuditRecord(ctx, arg)
	if err != nil {
		return nil, err
	}
	return toAuditRecordEntity(&model, ar), nil
}

func toAuditRecordEntity(model *AuditRecord, ent *entity.AuditRecord) *entity.AuditRecord {
	if ent == nil {
		ent = &entity.AuditRecord{}
	}

	ent.Action = audit.Action(model.Action)
	ent.CreatedAt = model.CreatedAt
	ent.Data = model.Data
	ent.ID = model.ID
	ent.OriginIP = model.OriginIp
	ent.ResourceID = model.ResourceID
	ent.ResourceType = audit.ResourceType(model.ResourceType)
	ent.UserID = model.UserID

	return ent
}
