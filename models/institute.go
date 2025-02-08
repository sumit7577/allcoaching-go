package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Institute represents the institute model.
type Institute struct {
	Id           int64     `orm:"auto"`
	Name         string    `orm:"size(100); null" valid:"MaxSize(100)"`
	About        string    `orm:"type(text); null"`
	Category     *Category `orm:"rel(fk); null"`
	Banner       []*Banner `orm:"rel(m2m); null"`
	DirectorName string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	User         *User     `orm:"rel(one); unique; notnull"`
	Image        string    `orm:"size(300); null"`
	DateCreated  time.Time `orm:"auto_now_add;type(datetime)"`
	DateUpdated  time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Institute))
}

func GetAllInstitues() (num int64, institutes []*Institute, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	num, err = qs.All(&institutes)
	if err != nil {
		return num, nil, err
	} else {
		return num, institutes, nil
	}
}

type CourseSerializer struct {
	Course    []*Course       `json:"course"`
	Videos    []*CourseVideos `json:"videos"`
	Institute *Institute      `json:"institute"`
}

func GetInstitute(uid int64, page int) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var courses []*Course
	var videos []*CourseVideos
	institute := Institute{}

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("course").Filter("Institute__Id", uid),
	}

	_, err := query.Paginate().All(&courses, "id", "name", "description", "price", "image", "category", "created_at", "updated_at")

	if err != nil {
		return nil, err
	}
	_, err = o.QueryTable("course_videos").
		Filter("Course__Institute__Id", uid).Offset(page).Limit(10).All(&videos)

	if err != nil {
		return nil, err
	}

	errs := o.QueryTable("institute").Filter("Id", uid).One(&institute)

	if errs != nil {
		return nil, errs
	}

	serializer := &CourseSerializer{
		Course:    courses,
		Videos:    videos,
		Institute: &institute,
	}

	data, err := query.CreatePagination(serializer)

	if err != nil {
		return nil, err
	}

	return data, nil

}

type InstituteSerializer struct {
	Category   *Category    `json:"category"`
	Institutes []*Institute `json:"institutes"`
}

func GetCategoriesWithInstitutes() ([]InstituteSerializer, error) {
	o := orm.NewOrm()
	var categories []Category
	_, err := o.QueryTable("category").All(&categories)
	if err != nil {
		return nil, err
	}

	var result []InstituteSerializer

	for _, category := range categories {
		var institutes []*Institute
		_, err := o.QueryTable("institute").Filter("Category__Id", category.Id).All(&institutes, "id", "name", "about", "director_name", "image", "date_created", "date_updated")
		if err != nil {
			return nil, err
		}
		result = append(result, InstituteSerializer{
			Category:   &category,
			Institutes: institutes,
		})
	}

	return result, nil
}
