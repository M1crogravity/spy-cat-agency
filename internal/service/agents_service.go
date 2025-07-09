package service

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
)

type AgentsRepository interface {
	FindByName(context.Context, string) (*model.Agent, error)
	Create(context.Context, *model.Agent) error
	FindById(context.Context, int64) (*model.Agent, error)
}

type AgentsService struct {
	repository AgentsRepository
}

func NewAgentsService(repo AgentsRepository) *AgentsService {
	return &AgentsService{
		repository: repo,
	}
}

func (s *AgentsService) GetByName(ctx context.Context, name string) (*model.Agent, error) {
	return s.repository.FindByName(ctx, name)
}

func (s *AgentsService) Create(ctx context.Context, agent *model.Agent) error {
	return s.repository.Create(ctx, agent)
}

func (s *AgentsService) GetById(ctx context.Context, id int64) (*model.Agent, error) {
	return s.repository.FindById(ctx, id)
}
