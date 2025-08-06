package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

type JWTIf interface {
	GenerateToken(userID uuid.UUID) (string, error)
	Decode(tokenString string, claims *Claims) error
}

type JWTStruct struct {
	ExpiredTime time.Duration
	SecretKey   string
}

var JWT = getJwt()

func getJwt() JWTIf {
	return &JWTStruct{
		ExpiredTime: env.AppEnv.JWTExpiredTime,
		SecretKey:   env.AppEnv.JWTSecretKey,
	}
}

func (j *JWTStruct) GenerateToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.AppEnv.AppName,
			Subject:   env.AppEnv.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
