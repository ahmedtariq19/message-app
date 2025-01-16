package services

import "message-app/models"

type AuthenticationService interface {
	CreateAuthToken(req *models.AuthenticationReq) (string, error)
}

type authenticationService struct {
	jwtService       JWTService
	gbeConfigService GbeConfigService
}

func NewAuthenticationService(
	jwtService JWTService,
	gbeConfigService GbeConfigService,
) AuthenticationService {
	return &authenticationService{
		jwtService:       jwtService,
		gbeConfigService: gbeConfigService,
	}
}

func (m *authenticationService) CreateAuthToken(req *models.AuthenticationReq) (string, error) {
	if m.gbeConfigService.GetConfig().Password != req.Password {
		return "", &models.StandardError{
			Code:        models.INVALID_PASSWORD,
			ActualError: nil,
			Line:        "CreateAuthToken():25",
			Message:     models.INVALID_PASSWORD_MSG,
		}
	}
	return m.jwtService.CreateLoginToken(req.Password)
}
