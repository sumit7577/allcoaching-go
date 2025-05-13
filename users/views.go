package users

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"allcoaching-go/utils"
	"fmt"
	"math/rand"
	"time"
)

type UserController struct {
	services.RestApi
}

// @Title Create User
// @Description Create a new user
// @Param	body	body	models.User	true	"User data in JSON format"
// @Success 200 {object} models.User
// @Failure 400 Bad Request
// @router / [post]
func (c *UserController) Post() {
	c.Models = &models.User{}
	c.Create(func() (interface{}, error) {
		manager := &models.UserManager{User: c.Models.(*models.User)}
		user, err := manager.CreateUser()
		if err != nil {
			message := utils.HandleUniqueConstraintError(err)
			return nil, message
		}
		return map[string]interface{}{
			"status":  "true",
			"message": "User created successfully",
			"data":    user,
		}, nil
	})
}

// @Title Get User
// @Description Get the current user
// @Success 200 {object} models.User
// @Failure 400 Bad Request
// @router / [get]
func (c *UserController) Get() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		if c.IsUserAuthenticated {
			return map[string]interface{}{
				"status": "true",
				"data":   c.CurrentUser,
			}, nil
		}

		return map[string]interface{}{
			"status":  "false",
			"message": "User Is Not Authenticated",
		}, nil
	})
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())               // always seed random
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 000000 to 999999
}

func (c *UserController) LoginUser() {
	c.Models = &models.LoginSerializer{}
	c.Create(func() (interface{}, error) {
		otp := generateOTP()
		phone := c.Models.(*models.LoginSerializer).Phone
		_, err := services.SendOtp(phone, otp)
		if err != nil {
			return nil, err
		}
		val, err := models.CreateOtp(phone, otp)

		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status":  "true",
			"message": "OTP sent successfully",
			"data":    val,
		}, nil
	})
}

// @Title Verify User
// @Description Verify the user with OTP
// @Param	body	body	models.OtpSerializer	true	"OTP data in JSON format"
// @Success 200 {object} models.User
// @Failure 400 Bad Request
// @router /verify [post]
func (c *UserController) VerifyUser() {
	c.Models = &models.OtpSerializer{}
	c.Create(func() (interface{}, error) {
		phone := c.Models.(*models.OtpSerializer).Phone
		otp := c.Models.(*models.OtpSerializer).Otp
		user, err := models.VerifyOtp(phone, otp)
		if err.Error() == "User not found" {
			return map[string]interface{}{
				"status":  "true",
				"message": "No User Found",
				"data":    user,
			}, nil
		}
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status":  "true",
			"message": "User verified successfully",
			"data":    user,
		}, nil
	})
}

func (c *UserController) CompleteUserVerify() {
	c.Models = &models.CompleteUserSignupSerializer{}
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.CompleteUserSignupSerializer)
			user := c.CurrentUser
			user.Name = request.Name
			user.Email = request.Email
			user.Username = models.GenerateRandomUsername(request.Name)
			user, err := models.CompleteUserVerification(c.CurrentUser)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"status":  "true",
				"message": "User verified successfully",
				"data":    user,
			}, nil
		})
		return nil, nil
	})
}
