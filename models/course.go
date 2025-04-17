package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

//Change Later and add digit to the price

type Course struct {
	Id          int64      `orm:"auto"`
	Name        string     `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Institute   *Institute `orm:"rel(fk); notnull"`
	Category    *Category  `orm:"rel(fk); null"`
	Banner      []*Banner  `orm:"rel(m2m); null"`
	Collection  string     `orm:"type(jsonb); null"`
	Description string     `orm:"type(text); null"`
	Price       float64    `orm:"decimal(2); notnull" valid:"Required"`
	Image       string     `orm:"size(300); null"`
	CreatedAt   time.Time  `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time  `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Course))
}
