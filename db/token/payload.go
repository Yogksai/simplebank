package token

import (
	"errors"
	"time"

	uuid "github.com/google/uuid"
)

var (
	ErrInvalidToken   = errors.New("token is invalid")
	ErrExpiredToken   = errors.New("token has expired")
	ErrInvalidKeySize = errors.New("key size must be at least 32 characters")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued"`
	ExpiredAt time.Time `json:"expired"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	if p.ID == uuid.Nil {
		return ErrInvalidToken
	}
	return nil
}
