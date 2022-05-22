package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

func (s *sqlStore) CreateStagedJob(ctx context.Context, sj *entity.StagedJob) (*entity.StagedJob, error) {
	_, err := s.db.NewInsert().
		Model(sj).
		Returning("*").
		Exec(ctx)

	return sj, err
}
