package models

import (
	"allcoaching-go/allcoachingProject"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Change Later and add digit to the price
type Course struct {
	Id          int64      `orm:"auto"`
	Name        string     `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	Institute   *Institute `orm:"rel(fk); notnull"`
	User        []*User    `orm:"rel(m2m); null"`
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

func GetCourse(user *User, page int, courseId int64) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var videos []*CourseVideos
	var testSeries []*TestSeries
	var documents []*Documents

	course := &Course{Id: courseId}
	if err := o.Read(course); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	var banners []*Banner
	if _, err := o.
		QueryTable("banner").
		Filter("Course__Course__Id", course.Id).
		All(&banners); err != nil {
		return nil, fmt.Errorf("failed to load banners: %w", err)
	}
	course.Banner = banners

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("course").Filter("Id", courseId),
	}

	if course.Price > 0 {
		user := &User{Id: user.Id}
		purchased := o.QueryM2M(course, "User").Exist(user)
		if !purchased {
			return nil, errors.New("course not purchased")
		}
	}

	// Fetch course data (both for free and paid if purchased)
	_, err := o.QueryTable("course_videos").
		Filter("Course__Id", course.Id).
		Limit(10).
		All(&videos)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %w", err)
	}

	_, err = o.QueryTable("test_series").
		Filter("Course__Id", course.Id).
		Limit(10).
		All(&testSeries, "id", "name", "description", "questions", "timer", "created_at", "updated_at")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test series: %w", err)
	}

	_, err = o.QueryTable("documents").
		Filter("Course__Id", course.Id).
		Limit(10).
		All(&documents, "id", "name", "description", "file", "created_at", "updated_at")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch documents: %w", err)
	}

	totalUser, err := o.QueryM2M(course, "User").Count()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch total users: %w", err)
	}

	serializer := &CourseSerializer{
		Course:         []*Course{course},
		Videos:         videos,
		Institute:      &Institute{Id: course.Institute.Id},
		TestSeries:     testSeries,
		Documents:      documents,
		Followed:       course.Price == 0 || o.QueryM2M(course, "User").Exist(&User{Id: user.Id}),
		FollowersCount: totalUser,
	}
	data, err := query.CreatePagination(serializer)

	if err != nil {
		return nil, err
	}
	return data, nil
}

type OrderResponse struct {
	Institute struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"institute"`
	Course struct {
		Name        string  `json:"name"`
		Image       string  `json:"image"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
	} `json:"course"`
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"user"`
	Status  string `json:"status"`
	OrderId string `json:"order_id"`
}

func firstNLines(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func PurchaseCourse(user *User, courseId int64) (*OrderResponse, error) {
	o := orm.NewOrm()
	course := &Course{Id: courseId}
	if err := o.Read(course); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	course.Institute = &Institute{Id: course.Institute.Id}

	if err := o.Read(course.Institute); err != nil {
		return nil, fmt.Errorf("institute not found: %w", err)
	}
	// Read user
	u := &User{Id: user.Id}
	if err := o.Read(u); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	//course users is m2m of User in beego the table name is course_users is created by beego orm automatically
	exists := o.QueryTable("course_users").Filter("user_id", user.Id).Filter("course_id", courseId).Exist()
	if exists {
		return nil, errors.New("course already purchased")
	}
	orderData := map[string]interface{}{
		"amount":   course.Price * 100, // Razorpay expects amount in paise
		"currency": "INR",
		"receipt":  fmt.Sprintf("receipt_user%d_course%d", user.Id, course.Id),
	}
	order, err := allcoachingProject.RazorpayClient.Order.Create(orderData, nil)
	if err != nil {
		return nil, fmt.Errorf("razorpay order creation failed: %w", err)
	}
	orderJson, _ := json.Marshal(order)

	// Save order to DB
	newOrder := &Order{
		Course:    course,
		User:      user,
		OrderData: string(orderJson),
		OrderId:   order["id"].(string),
		Status:    "PENDING",
	}
	if _, err := o.Insert(newOrder); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	resp := &OrderResponse{
		Status:  newOrder.Status,
		OrderId: newOrder.OrderId,
	}

	if newOrder.Course != nil {
		resp.Course.Name = newOrder.Course.Name
		resp.Course.Image = newOrder.Course.Image
		resp.Course.Price = newOrder.Course.Price
		resp.Course.Description = firstNLines(newOrder.Course.Description, 200)

		if course.Institute != nil {
			resp.Institute.Name = course.Institute.Name
			resp.Institute.Image = course.Institute.Image
		}
	}
	if newOrder.User != nil {
		resp.User.Name = newOrder.User.Name
		resp.User.Email = newOrder.User.Email
		resp.User.Phone = newOrder.User.Phone
	}

	return resp, nil
}
