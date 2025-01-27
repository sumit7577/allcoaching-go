package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Banner struct {
	Id          int64     `orm:"auto"`
	Title       string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Image       string    `orm:"size(200); notnull" valid:"Required; MaxSize(200)"`
	DateCreated time.Time `orm:"auto_now_add;type(datetime)"`
	DateUpdated time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Banner))
}
