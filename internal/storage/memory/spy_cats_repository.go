package memory

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type SpyCatsRepository struct {
	spyCats map[int64]*model.SpyCat
	names   map[string]int64
	lastId  int64
}

func NewSpyCatRepository() *SpyCatsRepository {
	return &SpyCatsRepository{
		spyCats: make(map[int64]*model.SpyCat),
		names:   make(map[string]int64),
	}
}

func (r *SpyCatsRepository) Create(ctx context.Context, spyCat *model.SpyCat) error {
	if _, ok := r.names[spyCat.Name]; ok {
		return storage.ErrorUniqueConstraintViolation
	}
	id := r.lastId + 1
	spyCat.Id = id
	r.spyCats[id] = spyCat
	r.names[spyCat.Name] = id
	r.lastId = id
	return nil
}

func (r *SpyCatsRepository) FindById(ctx context.Context, id int64) (*model.SpyCat, error) {
	spyCat, ok := r.spyCats[id]
	if !ok {
		return nil, storage.ErrorModelNotFound
	}

	return spyCat, nil
}

func (r *SpyCatsRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := r.spyCats[id]; !ok {
		return storage.ErrorModelNotFound
	}

	delete(r.spyCats, id)
	return nil
}

func (r *SpyCatsRepository) Save(ctx context.Context, spyCat model.SpyCat) error {
	spyCatToUpdate, ok := r.spyCats[spyCat.Id]
	if !ok {
		return storage.ErrorModelNotFound
	}

	spyCatToUpdate.Salary = spyCat.Salary

	return nil
}

func (r *SpyCatsRepository) FindAll(ctx context.Context) ([]*model.SpyCat, error) {
	spyCats := make([]*model.SpyCat, 0, len(r.spyCats))
	for _, spyCat := range r.spyCats {
		spyCats = append(spyCats, spyCat)
	}

	return spyCats, nil
}

func (r *SpyCatsRepository) FindByName(ctx context.Context, name string) (*model.SpyCat, error) {
	id, ok := r.names[name]
	if !ok {
		return nil, storage.ErrorModelNotFound
	}

	spyCat, ok := r.spyCats[id]
	if !ok {
		panic("spy cat repository data inconsistency")
	}

	return spyCat, nil
}
