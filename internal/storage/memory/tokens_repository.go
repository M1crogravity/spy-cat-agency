package memory

import (
	"context"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
)

type TokensRepository struct {
	tokens map[string]*model.Token
}

func NewTokensRepository() *TokensRepository {
	return &TokensRepository{
		tokens: make(map[string]*model.Token),
	}
}

func (r *TokensRepository) Create(ctx context.Context, token *model.Token) error {
	if _, ok := r.tokens[token.Plaintext]; ok {
		return storage.ErrorUniqueConstraintViolation
	}

	r.tokens[token.Plaintext] = token
	return nil
}

func (r *TokensRepository) FindByPlaintext(ctx context.Context, tokenPLaintext string, scope string) (*model.Token, error) {
	token, ok := r.tokens[tokenPLaintext]
	if !ok || token.Scope != scope {
		return nil, storage.ErrorModelNotFound
	}

	return token, nil
}
