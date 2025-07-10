package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/postgres/sqlc"
)

type AgentsRepository struct {
	queries *sqlc.Queries
}

func NewAgentsRepository(conn sqlc.DBTX) *AgentsRepository {
	return &AgentsRepository{
		queries: sqlc.New(conn),
	}
}

func (r *AgentsRepository) FindByName(ctx context.Context, name string) (*model.Agent, error) {
	agent, err := r.queries.FindAgentByName(ctx, name)
	if err != nil {
		return nil, storage.ErrorModelNotFound
	}

	return &model.Agent{
		Id:       agent.ID,
		Name:     agent.Name,
		Password: *model.NewPasswordFromHash(agent.PasswordHash),
	}, nil
}

func (r *AgentsRepository) Create(ctx context.Context, agent *model.Agent) error {
	id, err := r.queries.CreateAgent(ctx, sqlc.CreateAgentParams{
		Name:         agent.Name,
		PasswordHash: agent.Password.Hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return storage.ErrorUniqueConstraintViolation
		}
		return err
	}

	agent.Id = id

	return nil
}

func (r *AgentsRepository) FindById(ctx context.Context, id int64) (*model.Agent, error) {
	agent, err := r.queries.FindAgentById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Agent{
		Id:       agent.ID,
		Name:     agent.Name,
		Password: *model.NewPasswordFromHash(agent.PasswordHash),
	}, nil
}
