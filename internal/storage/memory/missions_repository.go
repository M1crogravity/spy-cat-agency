package memory

import (
	"context"
	"maps"
	"slices"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type MissionsRepository struct {
	missions      map[int64]*model.Mission
	lastMissionId int64
}

func NewMissionsRepository() *MissionsRepository {
	return &MissionsRepository{
		missions: make(map[int64]*model.Mission),
	}
}

func (r *MissionsRepository) CreateMission(ctx context.Context, mission *model.Mission) error {
	missionId := r.lastMissionId + 1
	r.lastMissionId = missionId
	mission.Id = missionId

	for _, target := range mission.Targets {
		err := r.CreateTarget(ctx, mission, target)
		if err != nil {
			return err
		}
	}
	r.missions[missionId] = mission
	return nil
}

func (r *MissionsRepository) FindMissionById(ctx context.Context, id int64) (*model.Mission, error) {
	mission, ok := r.missions[id]
	if !ok {
		return nil, storage.ErrorModelNotFound
	}
	return mission, nil
}

func (r *MissionsRepository) DeleteMission(ctx context.Context, id int64) error {
	if _, ok := r.missions[id]; !ok {
		return storage.ErrorModelNotFound
	}

	delete(r.missions, id)
	return nil
}

func (r *MissionsRepository) SaveMission(ctx context.Context, mission *model.Mission) error {
	if _, ok := r.missions[mission.Id]; !ok {
		return storage.ErrorModelNotFound
	}
	for _, target := range mission.Targets {
		if target.Id != 0 {
			target.MissionId = mission.Id
			r.SaveTarget(ctx, target)
			continue
		}
		r.CreateTarget(ctx, mission, target)
	}
	r.missions[mission.Id] = mission
	return nil
}

func (r *MissionsRepository) SaveTarget(ctx context.Context, target *model.Target) error {
	return nil
}

func (r *MissionsRepository) CreateTarget(ctx context.Context, mission *model.Mission, target *model.Target) error {
	target.MissionId = mission.Id
	target.Id = int64(len(mission.Targets) + 1)
	return nil
}

func (r *MissionsRepository) DeleteTarget(ctx context.Context, mission *model.Mission, targetId int64) error {
	var targetIndex int
	for i, missionTarget := range mission.Targets {
		if missionTarget.Id != targetId {
			continue
		}
		targetIndex = i
	}
	mission.Targets = append(mission.Targets[:targetIndex], mission.Targets[targetIndex+1:]...)
	return nil
}

func (r *MissionsRepository) FindActiveMission(ctx context.Context, spyCatId int64) (*model.Mission, error) {
	for _, mission := range r.missions {
		if mission.AssignedCatId == spyCatId && !mission.IsCompleted() {
			return mission, nil
		}
	}

	return nil, storage.ErrorModelNotFound
}

func (r *MissionsRepository) FindAll(ctx context.Context) ([]*model.Mission, error) {
	return slices.Collect(maps.Values(r.missions)), nil
}
