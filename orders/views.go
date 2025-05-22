package orders

import (
	"allcoaching-go/models"
	"allcoaching-go/services"
	"fmt"

	"github.com/razorpay/razorpay-go/utils"
)

type OrderController struct {
	services.RestApi
}

// @Title Get Order Events
// @Description Get order events
// @Param	uid		path 	string	true		"Order ID"
// @Success 200 {object} []models.OrderEvent
// @Failure 400 Bad Request
// change later signature
func (c *OrderController) GetOrderEvents() {
	c.Models = &models.OrderRazorPayWebhookEvent{}
	c.Create(func() (interface{}, error) {
		signature := c.Ctx.Input.Header("X-Razorpay-Signature")
		body := c.Ctx.Input.RequestBody
		status := utils.VerifyWebhookSignature(string(body), signature, "abcd")
		if !status {
			return nil, fmt.Errorf("signature verification failed: %v", status)
		}
		event := c.Models.(*models.OrderRazorPayWebhookEvent)
		switch event.Event {
		case "order.paid":
			err := models.CompleteOrder(signature, *event)
			if err != nil {
				return nil, fmt.Errorf("error completing order: %v", err)
			}
		default:
			fmt.Println("Unhandled event:", event.Event)
		}
		return map[string]interface{}{
			"status":  "true",
			"message": "Order events retrieved successfully",
		}, nil
	})
}
