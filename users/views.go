package users

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"allcoaching-go/utils"
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
