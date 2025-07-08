package service

import (
	"testing"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/memory"
)

func TestMissionsCreate(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						Name:    "target1",
						Country: "ua",
						State:   model.Created,
					},
					{
						Name:    "target2",
						Country: "us",
						State:   model.Created,
					},
				},
			},
			errCheck: nil,
		},
		{
			name: "no targets mission",
			mission: &model.Mission{
				Targets: []*model.Target{},
			},
			errCheck: ErrTooFewTargets,
		},
		{
			name: "too many targets",
			mission: &model.Mission{
				Targets: []*model.Target{
					{},
					{},
					{},
					{},
				},
			},
			errCheck: ErrTooMuchTargets,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			if tt.mission.Id == 0 {
				t.Fatal("mission id is not set")
			}
			for _, target := range tt.mission.Targets {
				if target.Id == 0 {
					t.Fatal("target id is not set")
				}
			}
		})
	}
}

func TestMissionsRemove(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		errCheck error
	}{
		{
			name: "unassigned mission",
			mission: &model.Mission{
				Targets: []*model.Target{
					{},
				},
			},
			errCheck: nil,
		}, {
			name: "assigned mission",
			mission: &model.Mission{
				AssignedCatId: 1,
				Targets: []*model.Target{
					{},
				},
			},
			errCheck: ErrCantDeleteMission,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}
			err = service.RemoveMission(t.Context(), tt.mission.Id)
			if err != tt.errCheck {
				t.Fatal(err)
			}
		})
	}
}

func TestCompleteMission(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		spyCat   *model.SpyCat
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				State: model.Created,
				Targets: []*model.Target{
					{
						State: model.Created,
					},
				},
			},
			errCheck: nil,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}
			mission, err := service.CompleteMission(t.Context(), tt.mission.Id)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			if !mission.IsCompleted() {
				t.Fatal("mission is not completed")
			}
		})
	}
}

func TestCompleteTarget(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		target   func(*model.Mission) *model.Target
		spyCat   *model.SpyCat
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.InProgress,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat:   &model.SpyCat{},
			errCheck: nil,
		},
		{
			name: "completed target",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.Completed,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat:   &model.SpyCat{},
			errCheck: nil,
		},
		{
			name: "target with no mission",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.InProgress,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return &model.Target{}
			},
			spyCat:   &model.SpyCat{},
			errCheck: storage.ErrorModelNotFound,
		},
		{
			name: "wrong spy cat",
			mission: &model.Mission{
				AssignedCatId: 1,
				Targets: []*model.Target{
					{
						State: model.InProgress,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat: &model.SpyCat{
				Id: 2,
			},
			errCheck: ErrAccessDenied,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}
			target := tt.target(tt.mission)
			err = service.CompleteTarget(t.Context(), tt.mission, target.Id, tt.spyCat)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			if !target.IsCompleted() {
				t.Fatal("target not completed")
			}
		})
	}
}

func TestUpdateTargetNotes(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		target   func(*model.Mission) *model.Target
		spyCat   *model.SpyCat
		notes    string
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				Id:            1,
				AssignedCatId: 1,
				State:         model.InProgress,
				Targets: []*model.Target{
					{
						State: model.InProgress,
						Notes: "notes",
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			notes:    "new notes",
			errCheck: nil,
		},
		{
			name: "complete target",
			mission: &model.Mission{
				Id:            1,
				AssignedCatId: 1,
				State:         model.InProgress,
				Targets: []*model.Target{
					{
						State: model.Completed,
						Notes: "notes",
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			notes:    "new notes",
			errCheck: ErrOperationNotAllowedOnCompleted,
		},
		{
			name: "unassigned target",
			mission: &model.Mission{
				Id:            1,
				AssignedCatId: 1,
				State:         model.InProgress,
				Targets: []*model.Target{
					{
						Notes: "notes",
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return &model.Target{
					Notes: "notes",
				}
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			notes:    "new notes",
			errCheck: storage.ErrorModelNotFound,
		},
		{
			name: "wrong spy cat",
			mission: &model.Mission{
				Id:            1,
				AssignedCatId: 1,
				State:         model.InProgress,
				Targets: []*model.Target{
					{
						State: model.InProgress,
						Notes: "notes",
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat: &model.SpyCat{
				Id: 2,
			},
			notes:    "new notes",
			errCheck: ErrAccessDenied,
		},
		{
			name: "complete mission",
			mission: &model.Mission{
				Id:            1,
				AssignedCatId: 1,
				State:         model.Completed,
				Targets: []*model.Target{
					{
						State: model.InProgress,
						Notes: "notes",
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			notes:    "new notes",
			errCheck: ErrOperationNotAllowedOnCompleted,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}
			target := tt.target(tt.mission)
			err = service.UpdateNotes(t.Context(), tt.mission, target.Id, tt.notes, tt.spyCat)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			mission, err := service.GetMissionByID(t.Context(), tt.mission.Id)
			if err != nil {
				t.Fatal(err)
			}
			target = tt.target(mission)
			if target.Notes != tt.notes {
				t.Fatal("notes haven't been updated")
			}
		})
	}
}

func TestRemoveTarget(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		target   func(*model.Mission) *model.Target
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.InProgress,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			errCheck: nil,
		},
		{
			name: "wrong target",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.InProgress,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return &model.Target{
					Id: -1,
				}
			},
			errCheck: ErrMissionTargetMissmatch,
		},
		{
			name: "completed target",
			mission: &model.Mission{
				Targets: []*model.Target{
					{
						State: model.Completed,
					},
				},
			},
			target: func(m *model.Mission) *model.Target {
				return m.Targets[0]
			},
			errCheck: ErrOperationNotAllowedOnCompleted,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)
			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}
			target := tt.target(tt.mission)
			err = service.RemoveTarget(t.Context(), tt.mission, target.Id)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			for _, missionTarget := range tt.mission.Targets {
				if missionTarget.Id == target.Id {
					t.Fatal("target was not removed")
				}
			}
		})
	}
}

func TestAddTarget(t *testing.T) {
	tc := []struct {
		name     string
		mission  *model.Mission
		target   *model.Target
		errCheck error
	}{
		{
			name: "happy path",
			mission: &model.Mission{
				Targets: []*model.Target{
					{},
				},
			},
			target:   &model.Target{},
			errCheck: nil,
		},
		{
			name: "mission completed",
			mission: &model.Mission{
				State: model.Completed,
				Targets: []*model.Target{
					{},
				},
			},
			target:   &model.Target{},
			errCheck: ErrOperationNotAllowedOnCompleted,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)

			err := service.CreateMission(t.Context(), tt.mission)
			if err != nil {
				t.Fatal(err)
			}

			err = service.AddTarget(t.Context(), tt.mission, tt.target)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			var found bool
			for _, target := range tt.mission.Targets {
				if target.Id == tt.target.Id {
					found = true
				}
			}
			if !found {
				t.Fatal("target was not added to the mission")
			}
		})
	}
}

func TestAssignMission(t *testing.T) {
	tc := []struct {
		name     string
		missions []*model.Mission
		mission  func([]*model.Mission) *model.Mission
		spyCat   *model.SpyCat
		errCheck error
	}{
		{
			name: "happy path",
			missions: []*model.Mission{
				{
					Targets: []*model.Target{
						{},
					},
				},
			},
			mission: func(m []*model.Mission) *model.Mission {
				return m[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			errCheck: nil,
		},
		{
			name: "already assigned to same spy cat",
			missions: []*model.Mission{
				{
					AssignedCatId: 1,
					State:         model.InProgress,
					Targets: []*model.Target{
						{},
					},
				},
			},
			mission: func(m []*model.Mission) *model.Mission {
				return m[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			errCheck: nil,
		},
		{
			name: "already assigned to different spy cat",
			missions: []*model.Mission{
				{
					AssignedCatId: 2,
					Targets: []*model.Target{
						{},
					},
				},
			},
			mission: func(m []*model.Mission) *model.Mission {
				return m[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			errCheck: ErrAlreadyAssigned,
		},
		{
			name: "spy cat is busy",
			missions: []*model.Mission{
				{
					Targets: []*model.Target{
						{},
					},
				},
				{
					AssignedCatId: 1,
					Targets: []*model.Target{
						{},
					},
				},
			},
			mission: func(m []*model.Mission) *model.Mission {
				return m[0]
			},
			spyCat: &model.SpyCat{
				Id: 1,
			},
			errCheck: ErrSpyCatIsBusy,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewMissionsRepository()
			service := NewMissionsService(repo)

			for _, mission := range tt.missions {
				err := service.CreateMission(t.Context(), mission)
				if err != nil {
					t.Fatal(err)
				}
			}

			mission := tt.mission(tt.missions)
			err := service.AssignMission(t.Context(), mission, tt.spyCat)
			if err != tt.errCheck {
				t.Fatal(err)
			}
			if err != nil {
				return
			}
			if mission.AssignedCatId != tt.spyCat.Id {
				t.Error("mission is not assigned")
			}
			if mission.State != model.InProgress {
				t.Error("mission is not in progress")
			}
		})
	}
}
