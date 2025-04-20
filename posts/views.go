package posts

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
)

type PostController struct {
	services.RestApi
}

func (c *PostController) Get() {
	//c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		page, _ := c.GetInt("page")
		posts, err := models.GetAllPosts(page)
		if err != nil {
			return nil, errors.New("Posts Not Found")
		}

		return map[string]interface{}{
			"status": "true",
			"data":   posts,
		}, nil
	})
}
