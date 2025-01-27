package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// User represents the user model.
type User struct {
	Id          uint64    `orm:"auto"`
	Name        string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Username    string    `orm:"size(100); null; unique"`
	Email       string    `orm:"size(100); notnull; unique" valid:"Required; MaxSize(100); Email"`
	Password    string    `orm:"size(400); notnull" valid:"Required"`
	Phone       string    `orm:"size(13); notnull; unique" valid:"Required; Match(^\\+\\d{2}\\d{10}$)"`
	State       string    `orm:"size(100); null" valid:"MaxSize(100)"`
	Pincode     uint64    `orm:"digit(10); null"`
	Address     string    `orm:"size(400); null" valid:"MaxSize(400)"`
	DateJoined  time.Time `orm:"auto_now_add;type(datetime)"`
	DateUpdated time.Time `orm:"auto_now;type(datetime)"`
	IsActive    bool      `orm:"default(true)"`
	Image       string    `orm:"size(300); null"`
	IsInstitute bool      `orm:"default(false)"`
}

type AuthToken struct {
	Key     string    `orm:"size(40);pk"`                 // Primary key
	Created time.Time `orm:"auto_now_add;type(datetime)"` // Auto set creation time
	User    *User     `orm:"rel(one); unique; notnull"`   // Foreign key to CustomUser
}

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(AuthToken))
}

func GetUserToken(token string) (*User, error) {
	o := orm.NewOrm()

	// Create an AuthToken object to load the related User
	authToken := AuthToken{}

	// Use QueryTable with RelatedSel to fetch the user in a single query
	err := o.QueryTable("auth_token").
		Filter("Key", token).
		RelatedSel().
		One(&authToken)
	if err != nil {
		return nil, err // Return error if token or user is not found
	}

	return authToken.User, nil
}
