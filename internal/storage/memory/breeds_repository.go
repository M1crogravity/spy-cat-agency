package memory

import "context"

type BreedsRepository struct {
	breeds []string
}

func NewBreedsRepository(breeds ...string) *BreedsRepository {
	return &BreedsRepository{
		breeds: breeds,
	}
}

func (r *BreedsRepository) FindAll(ctx context.Context) ([]string, error) {
	return r.breeds, nil
}
