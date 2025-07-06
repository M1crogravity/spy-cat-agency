package memory

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type MissionsRepository struct {
	missions      map[int64]*model.Mission
	targets       map[int64]*model.Target
	lastMissionId int64
	lastTargetId  int64
}

func NewMissionsRepository() *MissionsRepository {
	return &MissionsRepository{
		missions: make(map[int64]*model.Mission),
		targets:  make(map[int64]*model.Target),
	}
}

func (r *MissionsRepository) CreateMission(ctx context.Context, mission *model.Mission) error {
	missionId := r.lastMissionId + 1
	r.lastMissionId = missionId
	mission.Id = missionId

	for _, target := range mission.Targets {
		target.MissionId = missionId
		err := r.CreateTarget(ctx, target)
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
	for targetId, target := range r.targets {
		if target.MissionId != id {
			continue
		}
		delete(r.targets, targetId)
	}
	delete(r.missions, id)
	return nil
}

func (r *MissionsRepository) SaveMission(ctx context.Context, mission *model.Mission) error {
	if _, ok := r.missions[mission.Id]; !ok {
		return storage.ErrorModelNotFound
	}
	for _, target := range mission.Targets {
		target.MissionId = mission.Id
		if _, ok := r.targets[target.Id]; !ok {
			r.SaveTarget(ctx, target)
			continue
		}
		r.CreateTarget(ctx, target)
	}
	r.missions[mission.Id] = mission
	return nil
}

func (r *MissionsRepository) SaveTarget(ctx context.Context, target *model.Target) error {
	if _, ok := r.targets[target.Id]; !ok {
		return storage.ErrorModelNotFound
	}
	r.targets[target.Id] = target
	return nil
}

func (r *MissionsRepository) CreateTarget(ctx context.Context, target *model.Target) error {
	targetId := r.lastTargetId + 1
	r.lastTargetId = targetId
	target.Id = targetId
	r.targets[targetId] = target
	return nil
}

func (r *MissionsRepository) DeleteTarget(ctx context.Context, targetId int64) error {
	target, ok := r.targets[targetId]
	if !ok {
		return storage.ErrorModelNotFound
	}
	mission, ok := r.missions[target.MissionId]
	if !ok {
		return storage.ErrorModelNotFound
	}

	var targetIndex int
	for i, missionTarget := range mission.Targets {
		if missionTarget.Id != target.Id {
			continue
		}
		targetIndex = i
	}
	mission.Targets = append(mission.Targets[:targetIndex], mission.Targets[targetIndex+1:]...)
	delete(r.targets, target.Id)
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
