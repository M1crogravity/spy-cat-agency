package remote

import (
	"bytes"
	"context"
	"crypto/md5"
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
	cacheSum [16]byte
}

func NewBreedsRepository(client Client, cacheTtl time.Duration) *BreedsRepository {
	return &BreedsRepository{
		client:   client,
		cacheTtl: cacheTtl,
	}
}

func (r *BreedsRepository) FindAll(ctx context.Context) ([]string, error) {
	if r.cache == nil || time.Since(r.cachedAt) > r.cacheTtl {
		err := r.requestBreeds(ctx)
		if err != nil {
			return nil, err
		}
	}

	return r.cache, nil
}

func (r *BreedsRepository) requestBreeds(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", breedUrl, nil)
	if err != nil {
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Breed external API error")
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}
	hash := md5.Sum(buf.Bytes())
	if hash == r.cacheSum {
		r.cachedAt = time.Now()
		return nil
	}

	var breeds []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(&buf).Decode(&breeds); err != nil {
		return err
	}

	names := make([]string, len(breeds))
	for i, breed := range breeds {
		names[i] = breed.Name
	}

	r.cache = names
	r.cachedAt = time.Now()

	return nil
}
