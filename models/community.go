package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CommunityPost struct {
	Id        int64      `orm:"auto"`
	Institute *Institute `orm:"rel(fk); notnull"`
	Name      string     `orm:"size(500); notnull" valid:"Required; MaxSize(500)"`
	Options   string     `orm:"type(jsonb); null"`
	Image     string     `orm:"size(500); null" valid:"MaxSize(300)"`
	Link      string     `orm:"size(500); null" valid:"MaxSize(300)"`
	Type      string     `orm:"size(300); notnull" valid:"Required; MaxSize(300)"`
	CreatedAt time.Time  `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time  `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CommunityPost))
}

func GetAllPosts(page int) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var posts []*CommunityPost

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("community_post"),
	}
	_, errs := query.Paginate().RelatedSel("institute").All(&posts)
	if errs != nil {
		return nil, errs
	}
	data, err := query.CreatePagination(posts)
	if err != nil {
		return nil, err
	}
	return data, nil
}
