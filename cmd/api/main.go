package main

import (
	"context"
	"log"
	"time"

	_ "github.com/AmithSAI007/prj-flexdesk-user-service/docs"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/api"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/config"
	db "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/handler"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/middleware"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/service"
	"github.com/AmithSAI007/prj-flexdesk-user-service/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// @title           FlexDesk User Service API
// @version         1.0
// @description     This is the API for the FlexDesk User Service.

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 "Type 'Bearer' followed by a space and a JWT token."
func main() {

	cfg := config.LoadConfig()

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	dbConfig, err := pgxpool.ParseConfig(cfg.DBSource())
	if err != nil {
		logger.Fatal("Failed to parse database configuration", zap.Error(err))
	}
	dbConfig.MaxConns = int32(cfg.DBMaxConns)
	dbConfig.MinConns = int32(cfg.DBMinConns)
	dbConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	dbConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime
	dbConfig.HealthCheckPeriod = 5 * time.Minute // This can also be in config

	logger.Info("Connecting to database...",
		zap.String("db_user", cfg.DBUser),
		zap.String("db_name", cfg.DBName),
		zap.Bool("using_cloud_sql", cfg.CloudSqlConnectionName != ""),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		logger.Fatal("Failed to create database connection pool", zap.Error(err))
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connection established successfully")

	store := db.New(dbpool)

	keyService := token.NewLocalKeyService(logger)
	privateKey, publicKey, err := keyService.LoadKeys(ctx, cfg.PrivateKeyPath, cfg.PublicKeyPath)

	if err != nil {
		logger.Fatal("Failed to load cryptographic keys", zap.Error(err))
	}

	tokenService := service.NewTokenService(
		cfg.TokenIssuer,
		logger,
		privateKey,
		publicKey,
		store,
	)

	userService := service.NewUserService(logger, store)
	authService := service.NewAuthService(logger, userService, tokenService)

	authMiddleware := middleware.NewAuthMiddleware(logger, tokenService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.PrometheusMetrics())

	validate := validator.New()
	authHandler := handler.NewAuthHandler(logger, authService, validate)
	userHandler := handler.NewUserHandler(logger, userService, validate)
	handlers := &api.HandlerRegistry{
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleware,
		UserHandler:    userHandler,
	}
	api.SetupRoutes(router, handlers)

	logger.Info("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		logger.Error(err.Error())
	}
}
