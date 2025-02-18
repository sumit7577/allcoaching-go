package test_series

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
)

type TestSeriesController struct {
	services.RestApi
}

func (c *TestSeriesController) SubmitAnswer() {
	c.Models = &models.TestSeriesAttempt{}
	c.Permissions = []string{services.IsAuthenticated}

	c.ApiView(func() (interface{}, error) {
		c.Models.(*models.TestSeriesAttempt).User = c.CurrentUser

		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.TestSeriesAttempt)

			data, err := models.CreateAttempt(request)

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

func (c *TestSeriesController) Result() {
	c.Permissions = []string{services.IsAuthenticated}

	c.ApiView(func() (interface{}, error) {

		return nil, nil
	})

}

func (c *TestSeriesController) Submit() {
	c.Models = &models.SubmitTestSerializer{}
	c.Permissions = []string{services.IsAuthenticated}
	c.ApiView(func() (interface{}, error) {

		c.Create(func() (interface{}, error) {
			request := c.Models.(*models.SubmitTestSerializer)

			data, err := models.SubmitTestSeries(request, c.CurrentUser)

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
