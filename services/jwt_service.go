package services

import (
	"errors"
	"time"

	"message-app/models"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	VerifyAuthToken(tokenStr string) (string, error)
	CreateLoginToken(uid string) (string, error)
}
type jwtService struct {
	gbeConfig GbeConfigService
}

func NewJWTService(
	gbeConfig GbeConfigService,
) JWTService {
	return &jwtService{
		gbeConfig: gbeConfig,
	}
}

func (j *jwtService) VerifyAuthToken(tokenStr string) (string, error) {

	secret := j.gbeConfig.GetConfig().JwtSecret
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		stdErr := &models.StandardError{
			Code:        models.INVALID_TOKEN,
			ActualError: err,
			Line:        "VerifyToken():126",
			Message:     models.INVALID_TOKEN_MESSAGE,
		}
		return "", stdErr
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("cannot convert claim to MapClaims")
	}
	if !token.Valid {
		stdErr := &models.StandardError{
			Code:        models.INVALID_TOKEN,
			ActualError: err,
			Line:        "VerifyToken():126",
			Message:     models.INVALID_TOKEN_MESSAGE,
		}
		return "", stdErr
	}

	uidVal, found := claim["uid"]
	if !found {
		return "", errors.New("bad token")
	}
	uid := uidVal.(string)

	return uid, nil
}

func (j *jwtService) CreateLoginToken(uid string) (string, error) {

	claim := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := j.gbeConfig.GetConfig().JwtSecret

	return token.SignedString([]byte(secret))
}
