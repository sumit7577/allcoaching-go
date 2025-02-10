package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type TestSeries struct {
	Id          int64         `orm:"auto"`
	Course      *Course       `orm:"rel(fk); null"`
	Name        string        `orm:"size(300); null" valid:"MaxSize(300)"`
	Description string        `orm:"type(text); null"`
	File        string        `orm:"size(400); null"`
	Questions   string        `orm:"type(jsonb); notnull"`
	Timer       time.Duration `orm:"notnull"`
	CreatedAt   time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time     `orm:"auto_now;type(datetime)"`
}

type TestSeriesSolution struct {
	Id          int64       `orm:"auto"`
	TestSeries  *TestSeries `orm:"rel(fk); notnull"`
	Description string      `orm:"type(text); null"`
	Solution    string      `orm:"type(jsonb); notnull"`
	CreatedAt   time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time   `orm:"auto_now;type(datetime)"`
}

type TestSeriesAttempt struct {
	Id         int64       `orm:"auto"`
	TestSeries *TestSeries `orm:"rel(fk); notnull"`
	User       *User       `orm:"rel(fk); notnull"`
	Result     string      `orm:"type(jsonb); notnull"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(TestSeries))
	orm.RegisterModel(new(TestSeriesSolution))
	orm.RegisterModel(new(TestSeriesAttempt))
}
