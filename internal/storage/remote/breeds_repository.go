package remote

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const breedUrl = "https://api.thecatapi.com/v1/breeds"

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type BreedsRepository struct {
	client   Client
	cache    []string
	cacheTtl time.Duration
	cachedAt time.Time
}

func NewBreedsRepository(client Client, cacheTtl time.Duration) *BreedsRepository {
	return &BreedsRepository{
		client:   client,
		cacheTtl: cacheTtl,
	}
}

func (r *BreedsRepository) FindAll(ctx context.Context) ([]string, error) {
	if r.cache != nil && time.Since(r.cachedAt) < r.cacheTtl {
		return r.cache, nil
	}

	breeds, err := r.requestBreeds(ctx)
	if err != nil {
		return nil, err
	}

	r.cache = breeds
	r.cachedAt = time.Now()

	return r.cache, nil
}

func (r *BreedsRepository) requestBreeds(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", breedUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Breed external API error")
	}

	var breeds []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	names := make([]string, len(breeds))
	for i, breed := range breeds {
		names[i] = breed.Name
	}

	return names, nil
}
