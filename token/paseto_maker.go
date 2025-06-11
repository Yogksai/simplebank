package token

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type pasetoMaker struct {
	paseto       *paseto.V2
	symmetrucKey []byte
}

func NewPasetoMaker(symmetrucKey string) (Maker, error) {
	if len(symmetrucKey) < chacha20poly1305.KeySize {
		return nil, ErrInvalidKeySize
	}
	return &pasetoMaker{
		paseto:       paseto.NewV2(),
		symmetrucKey: []byte(symmetrucKey),
	}, nil
}

func (maker *pasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetrucKey, payload, nil)
}

// VerifyToken checks if the token is valid and returns the username if it is
func (maker *pasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	if err := maker.paseto.Decrypt(token, maker.symmetrucKey, payload, nil); err != nil {
		return nil, err
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
