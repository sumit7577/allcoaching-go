package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CourseLiveStream struct {
	Id          int64     `orm:"auto"`
	Name        string    `orm:"size(300); notnull" valid:"Required; MaxSize(300)"`
	Course      *Course   `orm:"rel(fk); notnull"`
	Description string    `orm:"type(text); null"`
	Live        string    `orm:"size(300); null"`
	Metadata    string    `orm:"type(jsonb); null"`
	Scheduled   time.Time `orm:"type(datetime); null"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseLiveStream))
}
