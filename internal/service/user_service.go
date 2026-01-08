package service

import (
	"context"
	"errors"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
)

type UserInterface interface {
	// Define user-related methods here, e.g., CreateUser, GetUser, UpdateUser, DeleteUser, etc.
	GetUserByEmail(ctx context.Context, q generated.Querier, email string) (*model.User, error)
	GetUserByID(ctx context.Context, q generated.Querier, userID uuid.UUID) (*model.User, error)
	CreateUser(ctx context.Context, q generated.Querier, username, email, password string) (*model.User, error)
}

type UserService struct {
	// Add necessary fields here, e.g., database connection, logger, etc.
	logger *zap.Logger
}

func NewUserService(logger *zap.Logger) UserInterface {
	return &UserService{
		logger: logger,
	}
}

var _ UserInterface = (*UserService)(nil)

func (s *UserService) GetUserByEmail(ctx context.Context, q generated.Querier, email string) (*model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database
	user, err := q.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		s.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, ErrInternalServer

	}

	return &model.User{
		ID:           user.ID.Bytes,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil

}

func (s *UserService) CreateUser(ctx context.Context, q generated.Querier, username, email, password string) (*model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database
	userParams := generated.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}

	createdUser, err := q.CreateUser(ctx, userParams)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, ErrInternalServer
	}

	return &model.User{
		ID:        createdUser.ID.Bytes,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, q generated.Querier, userID uuid.UUID) (*model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database

	user, err := q.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		s.logger.Error("Failed to get user by ID", zap.Error(err))
		return nil, ErrInternalServer
	}

	return &model.User{
		ID:           user.ID.Bytes,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil

}
