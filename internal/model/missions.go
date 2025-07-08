package model

type CompleteState string

const (
	Created    CompleteState = "created"
	InProgress CompleteState = "in_progress"
	Completed  CompleteState = "completed"
)

type Mission struct {
	Id            int64         `json:"id"`
	State         CompleteState `json:"state"`
	AssignedCatId int64         `json:"assigned_cat_id"`
	Targets       []*Target     `json:"targets"`
}

type Target struct {
	Id        int64         `json:"id"`
	MissionId int64         `json:"mission_id"`
	Name      string        `json:"name"`
	Country   string        `json:"country"`
	Notes     string        `json:"notes"`
	State     CompleteState `json:"state"`
}

func (m *Mission) Complete() {
	m.State = Completed
}

func (m *Mission) IsCompleted() bool {
	return m.State == Completed
}

func (m *Mission) IsAssignedToCat() bool {
	return m.AssignedCatId != 0
}

func (m *Mission) IsAssignedTo(sc *SpyCat) bool {
	return m.AssignedCatId == sc.Id
}

func (m *Mission) IsAllTargetsComplete() bool {
	for _, target := range m.Targets {
		if !target.IsCompleted() {
			return false
		}
	}

	return true
}

func (m *Mission) GetTarget(id int64) *Target {
	for _, target := range m.Targets {
		if target.Id == id {
			return target
		}
	}
	return nil
}

func (t *Target) IsCompleted() bool {
	return t.State == Completed
}

func (t *Target) Complete() {
	t.State = Completed
}

func (t *Target) UpdateNotes(notes string) {
	t.Notes = notes
}
