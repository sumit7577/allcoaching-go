package institute

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
	"fmt"
	"strconv"
)

type InstituteController struct {
	services.RestApi
}

func (c *InstituteController) Home() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		data, err := models.GetCategoriesWithInstitutes()
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status": "true",
			"data":   data,
		}, nil
	})
}

func (c *InstituteController) Get() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		id := c.GetString(":uid")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Institute ID")
			}

			ins, err := models.GetInstitute(uid)
			fmt.Printf("%+v\n", err)
			if err != nil {
				return nil, errors.New("Institute not found")
			}

			return map[string]interface{}{
				"status": "true",
				"data":   ins,
			}, nil

		}
		return nil, errors.New("Institute not found")
	})
}
