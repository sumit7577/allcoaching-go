package posts

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
	"strconv"
)

type PostController struct {
	services.RestApi
}

func (c *PostController) Get() {
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		page, _ := c.GetInt("page")
		posts, err := models.GetAllPosts(page, c.CurrentUser)
		if err != nil {
			return nil, errors.New("Posts Not Found")
		}

		return map[string]interface{}{
			"status": "true",
			"data":   posts,
		}, nil
	})
}

func (c *PostController) LikePost() {
	c.Permissions = []string{services.IsAuthenticated}
	c.Models = &models.LikePostSerializer{}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.LikePostSerializer)

			data, err := models.LikePost(request.PostID, c.CurrentUser)
			if err != nil {
				return nil, err

			}
			return map[string]interface{}{
				"status": "true",
				"data":   data,
			}, nil
		})
		return nil, nil
	})
}

func (c *PostController) GetAllPostComments() {
	//c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		id := c.GetString(":uid")
		page, _ := c.GetInt("page")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Post ID")
			}
			ins, err := models.GetAllPostComments(uid, page)
			if err != nil {
				return nil, errors.New("Post Comment Not found!")
			}

			return map[string]interface{}{
				"status": "true",
				"data":   ins,
			}, nil

		}
		return nil, errors.New("Post Comment Not found!")
	})
}

func (c *PostController) CreatePostComment() {
	c.Permissions = []string{services.IsAuthenticated}
	c.Models = &models.CommentPostSerializer{}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.CommentPostSerializer)

			data, err := models.CommentPost(request.PostID, c.CurrentUser, request.Comment)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"status": "true",
				"data":   data,
			}, nil
		})
		return nil, nil
	})
}
