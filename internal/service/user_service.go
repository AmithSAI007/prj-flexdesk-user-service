package service

import (
	"context"
	"errors"

	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/db"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/model"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/security"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
)

type UserInterface interface {
	// Define user-related methods here, e.g., CreateUser, GetUser, UpdateUser, DeleteUser, etc.
	RegisterUser(ctx context.Context, username, email, password string) (model.User, error)
}

type UserService struct {
	// Add necessary fields here, e.g., database connection, logger, etc.
	logger *zap.Logger
	store  db.Store
}

func NewUserService(logger *zap.Logger, store db.Store) UserInterface {
	return &UserService{
		logger: logger,
		store:  store,
	}
}

var _ UserInterface = (*UserService)(nil)

func (s *UserService) RegisterUser(ctx context.Context, username, email, password string) (model.User, error) {
	// Implement user registration logic here
	// For example, hash the password, validate input, and store the user in the database

	_, err := s.store.GetUserByEmail(ctx, email)
	if err == nil {
		s.logger.Warn("Attempt to register with existing email", zap.String("email", email))
		// Return an appropriate error indicating the email is already in use
		return model.User{}, ErrEmailAlreadyInUse
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("Failed to check existing email", zap.Error(err))
		return model.User{}, ErrInternalServer
	}

	hashedPassword, err := security.HashPassword(password) // Implement password hashing

	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return model.User{}, ErrInternalServer
	}

	userParams := generated.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.store.CreateUser(ctx, userParams)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return model.User{}, ErrInternalServer
	}

	return model.User{
		ID:        createdUser.ID.Bytes,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
	}, nil

}
