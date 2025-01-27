package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Category struct {
	Id          int64     `orm:"auto"`
	Name        string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Icon        string    `orm:"size(200); null"`
	DateCreated time.Time `orm:"auto_now_add;type(datetime)"`
	DateUpdated time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Category))
}
