package services

import (
	"allcoaching-go/models"
	"errors"
	"strings"
)

func (r *RestApi) IsAuthenticated() (bool, *models.User) {
	authHeader := r.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		return false, nil
	}
	if !strings.HasPrefix(authHeader, "Token ") {
		return false, nil
	}
	token := strings.TrimPrefix(authHeader, "Token ")
	user, err := models.GetUserToken(token)
	if err != nil {
		return false, nil
	}

	return true, user
}

func (r *RestApi) Authenticate() (*models.User, error) {
	value, user := r.IsAuthenticated()
	if value {
		return user, nil
	}

	return nil, errors.New("User is Not Authenticated")
}
