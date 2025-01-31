package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CourseVideo struct {
	Id          int64     `orm:"auto"`
	Name        string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Course      *Course   `orm:"rel(fk); notnull"`
	Category    *Category `orm:"rel(fk); null"`
	Description string    `orm:"type(text); null"`
	Video       string    `orm:"size(300); null"`
	Metadata    string    `orm:"type(json); null"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseVideo))
}
