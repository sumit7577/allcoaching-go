package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
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
