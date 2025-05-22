package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type Order struct {
	Id        int       `orm:"auto"`
	Course    *Course   `orm:"rel(fk); notnull"` // Foreign key to Course
	User      *User     `orm:"rel(fk); notnull"`
	OrderData string    `orm:"type(json);size(4096)"` // JSON as string
	Status    string    `orm:"size(100);default(PENDING)"`
	PaymentId string    `orm:"null;size(400)"`
	OrderId   string    `orm:"null;size(400)"`
	Signature string    `orm:"null;size(400)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Order))
}

type OrderRazorPayWebhookEvent struct {
	Entity    string                 `json:"entity"`
	AccountID string                 `json:"account_id"`
	Event     string                 `json:"event"`
	Contains  []string               `json:"contains"`
	Payload   map[string]interface{} `json:"payload"` // use map[string]interface{} for flexibility
	CreatedAt int64                  `json:"created_at"`
}

func CompleteOrder(signature string, event OrderRazorPayWebhookEvent) error {
	o := orm.NewOrm()
	paymentData, ok := event.Payload["payment"].(map[string]interface{})
	if !ok {
		return errors.New("invalid payment structure")
	}
	paymentEntity, ok := paymentData["entity"].(map[string]interface{})
	if !ok {
		return errors.New("missing payment entity")
	}
	orderID := paymentEntity["order_id"].(string)
	paymentID := paymentEntity["id"].(string)
	apiSecret, err := web.AppConfig.String("razorpay-api::apisecret")
	if err != nil {
		return fmt.Errorf("error reading API secret: %v", err)
	}
	params := map[string]interface{}{
		"razorpay_order_id":   orderID,
		"razorpay_payment_id": paymentID,
	}
	//success := utils.CheckSignature(signature, params, apiSecret)
	if !success {
		return fmt.Errorf("signature verification failed: %v", err)
	}
	order := Order{OrderId: orderID}
	err = o.Read(&order, "OrderId")
	if err != nil {
		return err
	}
	order.Status = "COMPLETED"
	order.PaymentId = event.Payload["payment_id"].(string)
	order.Signature = signature
	_, err = o.Update(&order)
	course := Course{Id: order.Course.Id}
	if err := o.Read(&course); err != nil {
		return err
	}

	m2m := o.QueryM2M(&course, "Users")
	if _, err := m2m.Add(order.User); err != nil {
		return fmt.Errorf("failed to add user to course: %v", err)
	}

	return err
}
