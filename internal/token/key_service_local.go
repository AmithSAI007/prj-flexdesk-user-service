package token

import (
	"context"
	"crypto/ecdsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type LocalKeyService struct {
	logger *zap.Logger
}

func NewLocalKeyService(logger *zap.Logger) KeyInterface {
	return &LocalKeyService{
		logger: logger,
	}
}

var _ KeyInterface = (*LocalKeyService)(nil)

func (lks *LocalKeyService) LoadKeys(ctx context.Context, privateKeyPath, publicKeyPath string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := lks.loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := lks.loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func (lks *LocalKeyService) loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	// Implementation to load private key from local file system
	data, err := os.ReadFile(path)
	if err != nil {
		lks.logger.Error("Failed to read private key file", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(data)
	if err != nil {
		lks.logger.Error("Failed to parse private key", zap.Error(err))
		return nil, err
	}

	return privateKey, nil
}

func (lks *LocalKeyService) loadPublicKey(path string) (*ecdsa.PublicKey, error) {
	// Implementation to load public key from local file system
	data, err := os.ReadFile(path)
	if err != nil {
		lks.logger.Error("Failed to read public key file", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM(data)
	if err != nil {
		lks.logger.Error("Failed to parse public key", zap.Error(err))
		return nil, err
	}

	return publicKey, nil
}
