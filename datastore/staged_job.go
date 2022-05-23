package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/uptrace/bun"
)

type StagedJob interface {
	FindAll(context.Context, ...repo.SelectCriteria) ([]entity.StagedJob, error)
	FindOne(context.Context, ...repo.SelectCriteria) (entity.StagedJob, error)
	Delete(context.Context, *entity.StagedJob) error
	Save(context.Context, *entity.StagedJob) error
	Update(context.Context, *entity.StagedJob) error
}

type stagedJobStore struct {
	*repo.CRUDRepository[entity.StagedJob]
}

func NewStagedJob(db bun.IDB) StagedJob {
	return stagedJobStore{&repo.CRUDRepository[entity.StagedJob]{DB: db}}
}
