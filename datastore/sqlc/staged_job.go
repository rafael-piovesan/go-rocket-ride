package sqlc

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/stagedjob"
)

func (s *sqlStore) CreateStagedJob(ctx context.Context, sj *entity.StagedJob) (*entity.StagedJob, error) {
	arg := CreateStagedJobParams{
		JobName: sj.JobName.String(),
		JobArgs: sj.JobArgs,
	}
	model, err := s.q.CreateStagedJob(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toStagedJobEntity(&model, sj), nil
}

func toStagedJobEntity(model *StagedJob, ent *entity.StagedJob) *entity.StagedJob {
	if ent == nil {
		ent = &entity.StagedJob{}
	}

	ent.ID = model.ID
	ent.JobName = stagedjob.JobName(model.JobName)
	ent.JobArgs = model.JobArgs

	return ent
}
