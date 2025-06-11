package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("failed to create payload: %w", err)
	}
	claims := jwt.RegisteredClaims{
		ID:        payload.ID.String(),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		Subject:   payload.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	return signedToken, err

}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	claims := &jwt.RegisteredClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, keyfunc)
	if err != nil {
		// Анализируем конкретную ошибку
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrExpiredToken
		default:
			return nil, ErrInvalidToken
		}
	}
	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}
	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid token ID: %w", err)
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  claims.Subject,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}
	return payload, nil
}
