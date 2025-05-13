package models

import (
	"errors"
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
	Phone       string    `orm:"size(10); notnull; unique" valid:"Required; Match(^\\d{10}$)"`
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
	Id      int64     `orm:"auto"`
	Key     string    `orm:"size(40);unique"`             // Primary key
	Created time.Time `orm:"auto_now_add;type(datetime)"` // Auto set creation time
	User    *User     `orm:"rel(one); unique; notnull"`   // Foreign key to CustomUser
}

type Otp struct {
	Id      int64     `orm:"auto"`
	Phone   string    `orm:"size(10); notnull; unique" valid:"Required; Match(^\\d{10}$)"`
	Otp     string    `orm:"size(6); notnull" valid:"Required; MaxSize(6)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(AuthToken))
	orm.RegisterModel(new(Otp))
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

type LoginSerializer struct {
	Phone string `json:"phone" valid:"Required; Match(^\\d{10}$)"`
}

func CreateOtp(phone string, otp string) (*Otp, error) {
	o := orm.NewOrm()
	otpModel := &Otp{
		Phone: phone,
		Otp:   otp,
	}
	_, err := o.Insert(otpModel)
	if err != nil {
		return nil, err
	}
	/*go func(phone string, otp string) {
		time.Sleep(30 * time.Second) // wait 30 seconds

		o := orm.NewOrm()
		_, err := o.QueryTable("otp").Filter("phone", phone).Filter("otp", otp).Delete()
		if err != nil {
			// Log the error, optional
			logs.Error("Failed to auto-delete OTP:", err)
		} else {
			logs.Info("OTP auto-deleted after 30 seconds for phone:", phone)
		}
	}(phone, otp)*/

	return otpModel, nil
}

type OtpSerializer struct {
	Phone string `json:"phone" valid:"Required; Match(^\\d{10}$)"`
	Otp   string `json:"otp" valid:"Required; MaxSize(6)"`
}

var ErrUserNotFound = errors.New("User not found")

func VerifyOtp(phone string, otp string) (*AuthToken, error) {
	o := orm.NewOrm()
	otpModel := &Otp{}
	err := o.QueryTable("otp").Filter("phone", phone).Filter("otp", otp).One(otpModel)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = o.QueryTable("user").Filter("phone", phone).One(user)
	if err == orm.ErrNoRows {
		newUser := &User{
			Phone: phone,
		}
		manager := &UserManager{User: newUser}
		user, err = manager.CreateUser()
		if err != nil {
			return nil, err
		}
		authToken := &AuthToken{}
		err = o.QueryTable("auth_token").Filter("user__id", user.Id).RelatedSel("user").One(authToken)
		if err == nil {
			return authToken, errors.New("User not found")
		}
		return nil, err
	}
	authToken := &AuthToken{}
	err = o.QueryTable("auth_token").Filter("user__id", user.Id).RelatedSel("user").One(authToken)
	if err == nil {
		return authToken, nil
	}
	return nil, err
}

type CompleteUserSignupSerializer struct {
	Name  string `json:"name" valid:"Required; MaxSize(150)"`
	Email string `json:"email" valid:"Required; MaxSize(100); Email"`
}

func CompleteUserVerification(user *User) (*User, error) {
	o := orm.NewOrm()
	_, err := o.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
