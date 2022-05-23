package datastore

import (
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/uptrace/bun"
)

type AuditRecord interface {
	data.ICRUDStore[entity.AuditRecord]
}

func NewAuditRecord(db bun.IDB) AuditRecord {
	return data.New[entity.AuditRecord](db)
}
