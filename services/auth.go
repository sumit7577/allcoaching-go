package services

import (
	"allcoaching-go/models"
	"errors"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/server/web"
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

func SendOtp(phone string, otp string) (string, error) {
	apiUrl, err := web.AppConfig.String("sms-api::apiurl")
	apiKey, err := web.AppConfig.String("sms-api::apikey")

	if err != nil {
		panic(err)
	}
	fullUrl := fmt.Sprintf("%s?authorization=%s&sender_id=%s&message=%s&variables_values=%s&route=dlt&numbers=%s", apiUrl, apiKey, "AllCOG", "184022", otp, phone)
	b := httplib.Get(fullUrl)

	str, err := b.String()
	if err != nil {
		panic(err)
	}

	return string(str), nil
}
