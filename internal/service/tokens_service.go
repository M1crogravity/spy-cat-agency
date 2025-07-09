package service

import (
	"context"
	"time"

	"github.com/m1crogravity/spy-cat-agency/internal/model"
)

type TokenRepository interface {
	Create(context.Context, *model.Token) error
	FindByPlaintext(context.Context, string, string) (*model.Token, error)
}

type TokenService struct {
	repository TokenRepository
}

func NewTokensService(repo TokenRepository) *TokenService {
	return &TokenService{
		repository: repo,
	}
}

func (s *TokenService) Create(ctx context.Context, userID int64, userType model.UserType, ttl time.Duration, scope string) (*model.Token, error) {
	token, err := model.GenerateToken(userID, userType, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, token)

	return token, err
}

func (s *TokenService) GetTokenByPlaintext(ctx context.Context, tokenPlaintext, tokenScope string) (*model.Token, error) {
	return s.repository.FindByPlaintext(ctx, tokenPlaintext, tokenScope)
}
