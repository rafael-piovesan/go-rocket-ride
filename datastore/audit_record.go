package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/uptrace/bun"
)

type AuditRecord interface {
	FindAll(context.Context, ...repo.SelectCriteria) ([]entity.AuditRecord, error)
	FindOne(context.Context, ...repo.SelectCriteria) (entity.AuditRecord, error)
	Delete(context.Context, *entity.AuditRecord) error
	Save(context.Context, *entity.AuditRecord) error
	Update(context.Context, *entity.AuditRecord) error
}

type auditRecordStore struct {
	*repo.CRUDRepository[entity.AuditRecord]
}

func NewAuditRecord(db bun.IDB) AuditRecord {
	return auditRecordStore{&repo.CRUDRepository[entity.AuditRecord]{DB: db}}
}
