package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1crogravity/spy-cat-agency/internal/model"
	"github.com/m1crogravity/spy-cat-agency/internal/storage"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/postgres/sqlc"
)

type TokensRepository struct {
	queries *sqlc.Queries
}

func NewTokensRepository(conn sqlc.DBTX) *TokensRepository {
	return &TokensRepository{
		queries: sqlc.New(conn),
	}
}

func (r *TokensRepository) Create(ctx context.Context, token *model.Token) error {
	return r.queries.CreateToken(ctx, sqlc.CreateTokenParams{
		Hash:     token.Hash,
		UserID:   token.UserID,
		UserType: string(token.UserType),
		Expiry:   pgtype.Timestamptz{Time: token.Expiry, Valid: true},
		Scope:    token.Scope,
	})
}

func (r *TokensRepository) FindByPlaintext(ctx context.Context, tokenPLaintext string, scope string) (*model.Token, error) {
	tokenHash := sha256.Sum256([]byte(tokenPLaintext))
	token, err := r.queries.FindTokenByPlaintext(ctx, sqlc.FindTokenByPlaintextParams{
		Hash:   tokenHash[:],
		Scope:  scope,
		Expiry: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, storage.ErrorModelNotFound
	}

	return &model.Token{
		Plaintext: base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(token.Hash),
		Hash:      token.Hash,
		UserID:    token.UserID,
		UserType:  model.UserType(token.UserType),
		Expiry:    token.Expiry.Time,
		Scope:     token.Scope,
	}, nil
}
