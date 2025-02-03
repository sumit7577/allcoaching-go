package services

import (
	"allcoaching-go/models"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
)

const (
	// authSessionKeyAuthenticated stores the key used to store the authentication status in the session
	IsAuthenticated = "authenticated"
	IsPagination    = "pagination"
)

// RestApi struct manages authentication, validation, and permissions
type RestApi struct {
	web.Controller
	validation.Validation
	Permissions         []string
	Models              interface{}
	IsUserAuthenticated bool
	CurrentUser         *models.User
}

// Helper function to check if a permission exists in the Permissions slice
func containsPermission(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (r *RestApi) ApiView(action func() (interface{}, error)) {

	if r.Permissions != nil {
		if containsPermission(r.Permissions, IsAuthenticated) {
			// Authenticate the user
			user, err := r.Authenticate()
			if err != nil {
				r.Data["json"] = map[string]string{"status": "false", "error": err.Error()}
				r.Ctx.Output.Status = http.StatusUnauthorized
				r.ServeJSON()
				return
			}
			r.CurrentUser = user
			r.IsUserAuthenticated = true
		}
	}

	data, err := action()
	if err != nil {
		r.Data["json"] = map[string]string{"status": "false", "error": err.Error()}
		r.Ctx.Output.Status = http.StatusBadRequest
		r.ServeJSON()
		return
	}

	r.Data["json"] = data
	r.ServeJSON()
}

func (r *RestApi) Create(action func() (interface{}, error)) {
	if r.Models != nil {
		json.Unmarshal(r.Ctx.Input.RequestBody, r.Models)
		valid, err := r.Valid(r.Models)

		if err != nil {
			r.Data["json"] = map[string]string{"status": "false", "error": "Validation error occurred: " + err.Error()}
			r.Ctx.Output.Status = http.StatusBadRequest
			r.ServeJSON()
			return
		}

		if !valid {
			errors := make(map[string]string)
			for _, err := range r.Errors {
				errors[err.Key] = err.Message
			}
			r.Data["json"] = map[string]interface{}{
				"error":  "Validation failed",
				"errors": errors,
				"status": "false",
			}
			r.Ctx.Output.Status = http.StatusBadRequest
			r.ServeJSON()
			return
		}

		if valid {
			data, err := action()
			if err != nil {
				r.Data["json"] = map[string]string{"status": "false", "error": err.Error()}
				r.Ctx.Output.Status = http.StatusBadRequest
				r.ServeJSON()
				return
			}
			r.Data["json"] = data
			r.ServeJSON()
		}
	}
}

func (r *RestApi) Get() {
	if r.Models != nil {
		json.Unmarshal(r.Ctx.Input.RequestBody, r.Models)
		valid, err := r.Valid(r.Models)

		if err != nil {
			r.Data["json"] = map[string]string{"status": "false", "error": "Validation error occurred: " + err.Error()}
			r.Ctx.Output.Status = http.StatusBadRequest
			r.ServeJSON()
			return
		}

		if !valid {
			errors := make(map[string]string)
			for _, err := range r.Errors {
				errors[err.Key] = err.Message
			}
			r.Data["json"] = map[string]interface{}{
				"error":  "Validation failed",
				"errors": errors,
				"status": "false",
			}
			r.Ctx.Output.Status = http.StatusBadRequest
			r.ServeJSON()
			return
		}
	}
}
