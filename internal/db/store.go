package db

import (
	"context"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store interface {
	// Add methods here that you want to expose from the store
	CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.FlexdeskUser, error)
	GetUserByEmail(ctx context.Context, email string) (generated.FlexdeskUser, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (generated.FlexdeskUser, error)
	GetUserByUsername(ctx context.Context, username string) (generated.FlexdeskUser, error)
	CreateRefreshToken(ctx context.Context, arg generated.CreateRefreshTokenParams) (generated.FlexdeskRefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (generated.FlexdeskRefreshToken, error)
	InvalidateRefreshToken(ctx context.Context, tokenId pgtype.UUID) error
	InvalidateRefreshTokenByHash(ctx context.Context, tokenHash string) error
}

var _ Store = (*generated.Queries)(nil)
