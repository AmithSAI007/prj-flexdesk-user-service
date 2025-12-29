package db

import (
	"context"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
)

type Store interface {
	// Add methods here that you want to expose from the store
	CreateUser(ctx context.Context, arg generated.CreateUserParams) (generated.FlexdeskUser, error)
	GetUserByEmail(ctx context.Context, email string) (generated.FlexdeskUser, error)
	GetUserByUsername(ctx context.Context, username string) (generated.FlexdeskUser, error)
}

var _ Store = (*generated.Queries)(nil)
