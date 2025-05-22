package courses

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
	"strconv"
)

type CourseController struct {
	services.RestApi
}

func (c *CourseController) Purchase() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		id := c.GetString(":uid")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Course ID")
			}
			data, err := models.PurchaseCourse(c.CurrentUser, uid)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"status": "true",
				"data":   data,
			}, nil
		}
		return nil, errors.New("Course not found")
	})
}
