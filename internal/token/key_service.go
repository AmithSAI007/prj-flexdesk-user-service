package token

import (
	"context"
	"crypto/ecdsa"
)

type KeyInterface interface {
	LoadKeys(ctx context.Context, privateKeyPath, publicKeyPath string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)
}
