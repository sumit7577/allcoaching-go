package orders

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"fmt"
)

type OrderController struct {
	services.RestApi
}

// @Title Get Order Events
// @Description Get order events
// @Param	uid		path 	string	true		"Order ID"
// @Success 200 {object} []models.OrderEvent
// @Failure 400 Bad Request

func (c *OrderController) GetOrderEvents() {
	c.Models = &models.OrderRazorPayWebhookEvent{}
	c.ApiView(func() (interface{}, error) {
		signature := c.Ctx.Input.Header("X-Razorpay-Signature")
		event := c.Models.(*models.OrderRazorPayWebhookEvent)
		switch event.Event {
		case "order.paid":
			err := models.CompleteOrder(signature, *event)
		default:
			fmt.Println("Unhandled event:", event.Event)
		}
		return map[string]interface{}{
			"status":  "true",
			"message": "Order events retrieved successfully",
		}, nil
	})
}
