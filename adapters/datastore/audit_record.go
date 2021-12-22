package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

func (s *sqlStore) CreateAuditRecord(ctx context.Context, ar *entity.AuditRecord) (*entity.AuditRecord, error) {
	_, err := s.db.NewInsert().
		Model(ar).
		Returning("*").
		Exec(ctx)

	return ar, err
}
