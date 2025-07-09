package service

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
)

type SpyCatsRepository interface {
	Create(context.Context, *model.SpyCat) error
	FindById(context.Context, int64) (*model.SpyCat, error)
	FindAll(context.Context) ([]*model.SpyCat, error)
	Save(context.Context, model.SpyCat) error
	Delete(context.Context, int64) error
	FindByName(context.Context, string) (*model.SpyCat, error)
}

type SpyCatService struct {
	repository SpyCatsRepository
}

func NewSpyCatService(repo SpyCatsRepository) *SpyCatService {
	return &SpyCatService{
		repository: repo,
	}
}

func (s *SpyCatService) Create(ctx context.Context, spyCat *model.SpyCat) error {
	return s.repository.Create(ctx, spyCat)
}

func (s *SpyCatService) UpdateSalary(ctx context.Context, id int64, salary float64) (*model.SpyCat, error) {
	spyCat, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	spyCat.Salary = salary

	return spyCat, s.repository.Save(ctx, *spyCat)
}

func (s *SpyCatService) GetAll(ctx context.Context) ([]*model.SpyCat, error) {
	return s.repository.FindAll(ctx)
}

func (s *SpyCatService) GetById(ctx context.Context, id int64) (*model.SpyCat, error) {
	return s.repository.FindById(ctx, id)
}

func (s *SpyCatService) Remove(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *SpyCatService) GetByName(ctx context.Context, name string) (*model.SpyCat, error) {
	return s.repository.FindByName(ctx, name)
}
