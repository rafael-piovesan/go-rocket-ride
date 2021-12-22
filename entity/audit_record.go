package entity

import (
	"encoding/json"
	"time"

	"github.com/rafael-piovesan/go-rocket-ride/entity/audit"
)

type AuditRecord struct {
	ID           int64
	Action       audit.Action
	CreatedAt    time.Time
	Data         json.RawMessage
	OriginIP     string
	ResourceID   int64
	ResourceType audit.ResourceType
	UserID       int64
}
