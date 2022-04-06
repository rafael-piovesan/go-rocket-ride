package jet

import (
	"context"
	"encoding/json"

	"github.com/rafael-piovesan/go-rocket-ride/entity"

	// dot import so go code would resemble as much as native SQL
	// dot import is not mandatory

	"github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/model"
	. "github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/table"
)

func (s *sqlStore) CreateAuditRecord(ctx context.Context, ar *entity.AuditRecord) (*entity.AuditRecord, error) {
	// RESEARCH: I think that JET can use entity classes that are compatible even if it didn't create them.
	// ANSWER: It seems you can't. According to: https://medium.com/@go.jet/jet-5f3667efa0cc
	// "developers have freedom to construct desired structure from auto-generated model files."

	data, err := json.Marshal(&ar.Data)
	if err != nil {
		panic(err)
	}

	newRecord := model.AuditRecords{
		ID:           ar.ID,
		Action:       ar.Action.String(),
		CreatedAt:    ar.CreatedAt,
		Data:         string(data),
		OriginIP:     ar.OriginIP,
		ResourceID:   ar.ResourceID,
		ResourceType: ar.ResourceType.String(),
		UserID:       ar.UserID,
	}

	insertStmt := AuditRecords.INSERT(AuditRecords.AllColumns).
		MODEL(newRecord).
		RETURNING(AuditRecords.AllColumns)

	dest := model.AuditRecords{}

	err = insertStmt.Query(s.db, &dest)

	return ar, err
}
