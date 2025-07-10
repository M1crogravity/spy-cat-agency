package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/postgres/sqlc"
)

type SpyCatsRepository struct {
	queries *sqlc.Queries
}

func NewSpyCatsRepository(conn sqlc.DBTX) *SpyCatsRepository {
	return &SpyCatsRepository{
		queries: sqlc.New(conn),
	}
}

func (r *SpyCatsRepository) Create(ctx context.Context, spyCat *model.SpyCat) error {
	id, err := r.queries.CreateSpyCat(ctx, sqlc.CreateSpyCatParams{
		Name:              spyCat.Name,
		PasswordHash:      spyCat.Password.Hash,
		YearsOfExperience: int32(spyCat.YearsOfExperience),
		Breed:             spyCat.Breed,
		Salary:            spyCat.Salary,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return storage.ErrorUniqueConstraintViolation
		}
	}
	spyCat.Id = id
	return nil
}

func (r *SpyCatsRepository) FindById(ctx context.Context, id int64) (*model.SpyCat, error) {
	spyCat, err := r.queries.FindSpyCatById(ctx, id)
	if err != nil {
		return nil, storage.ErrorModelNotFound
	}

	return convert(spyCat), nil
}

func (r *SpyCatsRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteSpyCatById(ctx, id)
	if err != nil {
		return storage.ErrorModelNotFound
	}
	return nil
}

func (r *SpyCatsRepository) Save(ctx context.Context, spyCat model.SpyCat) error {
	return r.queries.UpdateSpyCat(ctx, sqlc.UpdateSpyCatParams{
		ID:     spyCat.Id,
		Salary: spyCat.Salary,
	})
}

func (r *SpyCatsRepository) FindAll(ctx context.Context) ([]*model.SpyCat, error) {
	sc, err := r.queries.ListSpyCats(ctx)
	if err != nil {
		return nil, err
	}

	spyCats := make([]*model.SpyCat, len(sc))
	for i, spyCat := range sc {
		spyCats[i] = convert(spyCat)
	}
	return spyCats, nil
}

func (r *SpyCatsRepository) FindByName(ctx context.Context, name string) (*model.SpyCat, error) {
	spyCat, err := r.queries.FindSpyCatByName(ctx, name)
	if err != nil {
		return nil, storage.ErrorModelNotFound
	}

	return convert(spyCat), nil
}

func convert(spyCat sqlc.SpyCat) *model.SpyCat {
	return &model.SpyCat{
		Id:                spyCat.ID,
		Name:              spyCat.Name,
		Password:          model.Password{Hash: spyCat.PasswordHash},
		YearsOfExperience: int(spyCat.YearsOfExperience),
		Breed:             spyCat.Breed,
		Salary:            spyCat.Salary,
	}
}
