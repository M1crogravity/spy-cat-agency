package memory

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type AgentsRepository struct {
	agents map[int64]*model.Agent
	names  map[string]int64
	lastId int64
}

func NewAgentsRepository() *AgentsRepository {
	return &AgentsRepository{
		agents: make(map[int64]*model.Agent),
		names:  make(map[string]int64),
	}
}

func (r *AgentsRepository) FindByName(ctx context.Context, name string) (*model.Agent, error) {
	id, ok := r.names[name]
	if !ok {
		return nil, storage.ErrorModelNotFound
	}
	agent, ok := r.agents[id]
	if !ok {
		panic("agents repository data inconsistency")
	}
	return agent, nil
}

func (r *AgentsRepository) Create(ctx context.Context, agent *model.Agent) error {
	if _, ok := r.names[agent.Name]; ok {
		return storage.ErrorUniqueConstraintViolation
	}
	id := r.lastId + 1
	r.lastId = id
	agent.Id = id
	r.agents[id] = agent
	r.names[agent.Name] = id
	return nil
}

func (r *AgentsRepository) FindById(ctx context.Context, id int64) (*model.Agent, error) {
	agent, ok := r.agents[id]
	if !ok {
		return nil, storage.ErrorModelNotFound
	}
	return agent, nil
}
