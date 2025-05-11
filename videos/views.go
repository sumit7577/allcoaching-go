package videos

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"errors"
	"strconv"
)

type VideoController struct {
	services.RestApi
}

func (c *VideoController) GetAllVideoComments() {
	//c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {
		id := c.GetString(":uid")
		page, _ := c.GetInt("page")
		if id != "" {
			uid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, errors.New("Invalid Video ID")
			}
			ins, err := models.GetAllVideoComments(uid, page)
			if err != nil {
				return nil, errors.New("Video Comment Not found!")
			}

			return map[string]interface{}{
				"status": "true",
				"data":   ins,
			}, nil

		}
		return nil, errors.New("Post Comment Not found!")
	})
}

func (c *VideoController) CreateVideoComment() {
	c.Permissions = []string{services.IsAuthenticated}
	c.Models = &models.CommentVideoSerializer{}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.CommentVideoSerializer)

			data, err := models.CommentVideo(request.VideoID, c.CurrentUser, request.Comment)
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

func (c *VideoController) LikeVideo() {
	c.Permissions = []string{services.IsAuthenticated}
	c.Models = &models.LikeVideoSerializer{}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.LikeVideoSerializer)

			data, err := models.LikeVideo(request.VideoID, c.CurrentUser)
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

func (c *VideoController) IncreaseViewCount() {
	c.Permissions = []string{services.IsAuthenticated}
	c.Models = &models.LikeVideoSerializer{}
	c.ApiView(func() (interface{}, error) {
		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.LikeVideoSerializer)

			data, err := models.IncreaseViewCount(request.VideoID)
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
