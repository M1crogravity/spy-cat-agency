package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

const (
	MinTargets = 1
	MaxTargets = 3
)

var (
	ErrCantDeleteMission              = errors.New("can't delete the mission already assigned to the cat")
	ErrTooFewTargets                  = fmt.Errorf("can't create the mission. It should contain at least %d targets", MinTargets)
	ErrTooMuchTargets                 = fmt.Errorf("can't create the mission. It should contain maximum %d targets", MinTargets)
	ErrAccessDenied                   = errors.New("access denied")
	ErrNoMissionTarget                = errors.New("the target has no mission")
	ErrMissionTargetMissmatch         = errors.New("the mission doesn't have particular target")
	ErrOperationNotAllowedOnCompleted = errors.New("the operation is not allowed on completed subject")
	ErrAlreadyAssigned                = errors.New("the mission is already assigned")
	ErrSpyCatIsBusy                   = errors.New("can't assign the mission to busy spy cat")
)

type MissionsRepository interface {
	CreateMission(context.Context, *model.Mission) error
	CreateTarget(context.Context, *model.Target) error
	FindMissionById(context.Context, int64) (*model.Mission, error)
	DeleteMission(context.Context, int64) error
	SaveMission(context.Context, *model.Mission) error
	SaveTarget(context.Context, *model.Target) error
	DeleteTarget(context.Context, int64) error
	FindActiveMission(context.Context, int64) (*model.Mission, error)
	FindAll(context.Context) ([]*model.Mission, error)
}

type MissionsService struct {
	repository MissionsRepository
}

func NewMissionsService(repo MissionsRepository) *MissionsService {
	return &MissionsService{
		repository: repo,
	}
}

func (s *MissionsService) CreateMission(ctx context.Context, mission *model.Mission) error {
	targetsCount := len(mission.Targets)
	if targetsCount < MinTargets {
		return ErrTooFewTargets
	}
	if targetsCount > MaxTargets {
		return ErrTooMuchTargets
	}

	return s.repository.CreateMission(ctx, mission)
}

func (s *MissionsService) GetMissionByID(ctx context.Context, id int64) (*model.Mission, error) {
	return s.repository.FindMissionById(ctx, id)
}

func (s *MissionsService) RemoveMission(ctx context.Context, id int64) error {
	mission, err := s.GetMissionByID(ctx, id)
	if err != nil {
		return err
	}
	if mission.IsAssignedToCat() {
		return ErrCantDeleteMission
	}

	return s.repository.DeleteMission(ctx, id)
}

func (s *MissionsService) CompleteMission(ctx context.Context, id int64) (*model.Mission, error) {
	mission, err := s.repository.FindMissionById(ctx, id)
	if err != nil {
		return nil, err
	}

	if mission.IsCompleted() {
		return mission, nil
	}

	mission.Complete()

	return mission, s.repository.SaveMission(ctx, mission)
}

func (s *MissionsService) CompleteTarget(ctx context.Context, mission *model.Mission, targetId int64, spyCat *model.SpyCat) error {
	if !mission.IsAssignedTo(spyCat) {
		return ErrAccessDenied
	}

	target := mission.GetTarget(targetId)
	if target == nil {
		return storage.ErrorModelNotFound
	}

	if target.IsCompleted() {
		return nil
	}

	target.Complete()

	err := s.repository.SaveTarget(ctx, target)
	if err != nil {
		return err
	}

	if mission.IsAllTargetsComplete() {
		_, err = s.CompleteMission(ctx, mission.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MissionsService) UpdateNotes(ctx context.Context, mission *model.Mission, targetId int64, notes string, spyCat *model.SpyCat) error {
	if !mission.IsAssignedTo(spyCat) {
		return ErrAccessDenied
	}
	if mission.IsCompleted() {
		return ErrOperationNotAllowedOnCompleted
	}

	target := mission.GetTarget(targetId)
	if target == nil {
		return storage.ErrorModelNotFound
	}

	if target.IsCompleted() {
		return ErrOperationNotAllowedOnCompleted
	}

	target.UpdateNotes(notes)
	return s.repository.SaveTarget(ctx, target)
}

func (s *MissionsService) RemoveTarget(ctx context.Context, mission *model.Mission, targetId int64) error {
	found := mission.GetTarget(targetId)
	if found == nil {
		return ErrMissionTargetMissmatch
	}
	if found.IsCompleted() {
		return ErrOperationNotAllowedOnCompleted
	}

	return s.repository.DeleteTarget(ctx, targetId)
}

func (s *MissionsService) AddTarget(ctx context.Context, mission *model.Mission, target *model.Target) error {
	if mission.IsCompleted() {
		return ErrOperationNotAllowedOnCompleted
	}

	target.MissionId = mission.Id
	target.State = mission.State
	mission.Targets = append(mission.Targets, target)
	return s.repository.CreateTarget(ctx, target)
}

func (s *MissionsService) AssignMission(ctx context.Context, mission *model.Mission, spyCat *model.SpyCat) error {
	if mission.IsAssignedTo(spyCat) {
		return nil
	}
	if mission.IsAssignedToCat() {
		return ErrAlreadyAssigned
	}
	_, err := s.repository.FindActiveMission(ctx, spyCat.Id)
	switch {
	case errors.Is(err, storage.ErrorModelNotFound):
	case err == nil:
		return ErrSpyCatIsBusy
	default:
		return err
	}

	mission.AssignedCatId = spyCat.Id
	mission.State = model.InProgress
	return s.repository.SaveMission(ctx, mission)
}

func (s *MissionsService) GetAll(ctx context.Context) ([]*model.Mission, error) {
	return s.repository.FindAll(ctx)
}
