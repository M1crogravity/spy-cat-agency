package model

type CompleteState string

const (
	Created    CompleteState = "created"
	InProgress CompleteState = "in_progress"
	Completed  CompleteState = "completed"
)

type Mission struct {
	Id            int64
	State         CompleteState
	AssignedCatId int64
	Targets       []*Target
}

type Target struct {
	Id        int64
	MissionId int64
	Name      string
	Country   string
	Notes     string
	State     CompleteState
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

func (t *Target) IsCompleted() bool {
	return t.State == Completed
}

func (t *Target) Complete() {
	t.State = Completed
}

func (t *Target) UpdateNotes(notes string) {
	t.Notes = notes
}
