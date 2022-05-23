package datastore

import (
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/uptrace/bun"
)

type StagedJob interface {
	data.ICRUDStore[entity.StagedJob]
}

func NewStagedJob(db bun.IDB) StagedJob {
	return data.New[entity.StagedJob](db)
}
