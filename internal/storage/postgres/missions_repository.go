package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/postgres/sqlc"
)

type Connection interface {
	sqlc.DBTX
	Begin(context.Context) (pgx.Tx, error)
}

type MissionsRepository struct {
	queries    *sqlc.Queries
	connection Connection
}

func NewMissionsRepository(conn Connection) *MissionsRepository {
	return &MissionsRepository{
		queries:    sqlc.New(conn),
		connection: conn,
	}
}

func (r *MissionsRepository) CreateMission(ctx context.Context, mission *model.Mission) error {
	tx, err := r.connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txQuery := r.queries.WithTx(tx)
	id, err := txQuery.CreateMission(ctx, string(mission.State))
	if err != nil {
		return err
	}

	mission.Id = id

	arg := make([]sqlc.CreateTargetsParams, len(mission.Targets))
	for i, target := range mission.Targets {
		target.Id = int64(i + 1)
		target.MissionId = mission.Id
		arg[i] = sqlc.CreateTargetsParams{
			ID:        target.Id,
			MissionID: target.MissionId,
			Name:      target.Name,
			Country:   target.Country,
			Notes:     target.Notes,
			State:     string(target.State),
		}
	}

	_, err = txQuery.CreateTargets(ctx, arg)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *MissionsRepository) FindMissionById(ctx context.Context, id int64) (*model.Mission, error) {
	missionRows, err := r.queries.FindMissionById(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(missionRows) == 0 {
		return nil, storage.ErrorModelNotFound
	}

	targets := make([]*model.Target, len(missionRows))
	for i, missionRow := range missionRows {
		targets[i] = &model.Target{
			Id:        missionRow.TargetID,
			MissionId: missionRow.MissionID,
			Name:      missionRow.Name,
			Country:   missionRow.Country,
			Notes:     missionRow.Notes,
			State:     model.CompleteState(missionRow.TargetState),
		}
	}
	row := missionRows[0]

	return &model.Mission{
		Id:            row.MissionID,
		State:         model.CompleteState(row.MissionState),
		AssignedCatId: row.SpyCatID.Int64,
		Targets:       targets,
	}, nil
}

func (r *MissionsRepository) DeleteMission(ctx context.Context, id int64) error {
	err := r.queries.DeleteMission(ctx, id)
	if err != nil {
		return storage.ErrorModelNotFound
	}
	return nil
}

func (r *MissionsRepository) SaveMission(ctx context.Context, mission *model.Mission) error {
	tx, err := r.connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txQuery := r.queries.WithTx(tx)
	err = txQuery.UpdateMission(ctx, sqlc.UpdateMissionParams{
		ID:       mission.Id,
		State:    string(mission.State),
		SpyCatID: pgtype.Int8{Int64: mission.AssignedCatId, Valid: mission.AssignedCatId != 0},
	})
	if err != nil {
		return err
	}

	insertTargets := []sqlc.CreateTargetsParams{}
	for i, target := range mission.Targets {
		if target.Id == 0 {
			target.Id = int64(i + 1)
			target.MissionId = mission.Id
			insertTargets = append(insertTargets, sqlc.CreateTargetsParams{
				ID:        target.Id,
				MissionID: target.MissionId,
				Name:      target.Name,
				Country:   target.Country,
				Notes:     target.Notes,
				State:     string(target.State),
			})

			continue
		}

		err = updateTarget(ctx, txQuery, target)
		if err != nil {
			return err
		}
	}

	_, err = txQuery.CreateTargets(ctx, insertTargets)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *MissionsRepository) SaveTarget(ctx context.Context, target *model.Target) error {
	return updateTarget(ctx, r.queries, target)
}

func (r *MissionsRepository) CreateTarget(ctx context.Context, mission *model.Mission, target *model.Target) error {
	target.MissionId = mission.Id
	target.Id = int64(len(mission.Targets) + 1)
	_, err := r.queries.CreateTargets(ctx, []sqlc.CreateTargetsParams{
		{
			ID:        target.Id,
			MissionID: target.MissionId,
			Name:      target.Name,
			Country:   target.Country,
			Notes:     target.Notes,
			State:     string(target.State),
		},
	})

	return err
}

func (r *MissionsRepository) DeleteTarget(ctx context.Context, mission *model.Mission, targetId int64) error {
	err := r.queries.DeleteTarget(ctx, sqlc.DeleteTargetParams{
		ID:        targetId,
		MissionID: mission.Id,
	})
	if err != nil {
		return err
	}

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
	missionRows, err := r.queries.FindActiveMission(ctx, pgtype.Int8{Int64: spyCatId, Valid: true})
	if err != nil {
		return nil, err
	}
	if len(missionRows) == 0 {
		return nil, storage.ErrorModelNotFound
	}

	targets := make([]*model.Target, len(missionRows))
	for i, missionRow := range missionRows {
		targets[i] = &model.Target{
			Id:        missionRow.TargetID,
			MissionId: missionRow.MissionID,
			Name:      missionRow.Name,
			Country:   missionRow.Country,
			Notes:     missionRow.Notes,
			State:     model.CompleteState(missionRow.TargetState),
		}
	}
	row := missionRows[0]

	return &model.Mission{
		Id:            row.MissionID,
		State:         model.CompleteState(row.MissionState),
		AssignedCatId: row.SpyCatID.Int64,
		Targets:       targets,
	}, nil
}

func (r *MissionsRepository) FindAll(ctx context.Context) ([]*model.Mission, error) {
	missionRows, err := r.queries.FindAllMissions(ctx)
	if err != nil {
		return nil, err
	}

	var currentMission *model.Mission
	missions := []*model.Mission{}

	for _, missionRow := range missionRows {
		if currentMission == nil || currentMission.Id != missionRow.MissionID {
			currentMission = &model.Mission{
				Id:            missionRow.MissionID,
				State:         model.CompleteState(missionRow.MissionState),
				AssignedCatId: missionRow.SpyCatID.Int64,
				Targets:       make([]*model.Target, 0),
			}
			missions = append(missions, currentMission)
		}
		currentMission.Targets = append(currentMission.Targets, &model.Target{
			Id:        missionRow.TargetID,
			MissionId: missionRow.MissionID,
			Name:      missionRow.Name,
			Country:   missionRow.Country,
			Notes:     missionRow.Notes,
			State:     model.CompleteState(missionRow.TargetState),
		})
	}

	return missions, nil
}

func updateTarget(ctx context.Context, queries *sqlc.Queries, target *model.Target) error {
	return queries.UpdateTarget(ctx, sqlc.UpdateTargetParams{
		ID:        target.Id,
		MissionID: target.MissionId,
		Notes:     target.Notes,
		State:     string(target.State),
	})
}
