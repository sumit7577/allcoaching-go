package models

import "time"

type CourseVideos struct {
	Id          int64     `orm:"auto"`
	Name        string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Course      *Course   `orm:"rel(fk); notnull"`
	Category    *Category `orm:"rel(fk); null"`
	Description string    `orm:"type(text); null"`
	Image       string    `orm:"size(300); null"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}
