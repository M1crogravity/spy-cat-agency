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
}

type SpyCatService struct {
	repo SpyCatsRepository
}

func NewSpyCatService(repo SpyCatsRepository) *SpyCatService {
	return &SpyCatService{
		repo: repo,
	}
}

func (s *SpyCatService) Create(ctx context.Context, spyCat *model.SpyCat) error {
	err := s.repo.Create(ctx, spyCat)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpyCatService) UpdateSalary(ctx context.Context, id int64, salary float64) error {
	spyCat, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	spyCat.Salary = salary

	return s.repo.Save(ctx, *spyCat)
}

func (s *SpyCatService) GetAll(ctx context.Context) ([]*model.SpyCat, error) {
	return s.repo.FindAll(ctx)
}

func (s *SpyCatService) GetById(ctx context.Context, id int64) (*model.SpyCat, error) {
	return s.repo.FindById(ctx, id)
}

func (s *SpyCatService) Remove(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
