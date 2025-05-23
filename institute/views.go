package institute

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
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

func (c *InstituteController) GetAllHomeBanner() {
	//c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		data, err := models.GetAllHomeBanner()
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
		page, _ := c.GetInt("page")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Institute ID")
			}
			ins, err := models.GetInstitute(uid, page, c.CurrentUser)
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

func (c *InstituteController) GetAllCategories() {
	//c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		data, err := models.GetAllCategories()
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status": "true",
			"data":   data,
		}, nil
	})
}

func (c *InstituteController) GetSearchInstitute() {
	c.ApiView(func() (interface{}, error) {
		query := c.GetString("query")
		if query == "" {
			return nil, errors.New("No Institute found")
		}
		data, err := models.GetSearchIns(query)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"status": "true",
			"data":   data,
		}, nil
	})
}

func (c *InstituteController) ToggleFollow() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		id := c.GetString(":uid")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Institute ID")
			}
			message, err := models.ToggleFollowInstitute(uid, c.CurrentUser)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"status":  "true",
				"message": message,
			}, nil
		}
		return nil, errors.New("Institute not found")
	})
}
