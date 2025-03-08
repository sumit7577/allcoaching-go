package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Documents struct {
	Id          int64     `orm:"auto"`
	Course      *Course   `orm:"rel(fk); null"`
	Name        string    `orm:"size(300); null" valid:"MaxSize(300)"`
	Description string    `orm:"type(text); null"`
	File        string    `orm:"size(400); null"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Documents))
}
