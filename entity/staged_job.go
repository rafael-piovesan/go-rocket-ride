package entity

import (
	"encoding/json"

	"github.com/rafael-piovesan/go-rocket-ride/entity/stagedjob"
)

type StagedJob struct {
	ID      int64
	JobName stagedjob.JobName
	JobArgs json.RawMessage
}
