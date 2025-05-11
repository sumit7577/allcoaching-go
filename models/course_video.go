package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type CourseVideos struct {
	Id          int64     `orm:"auto"`
	Name        string    `orm:"size(300); notnull" valid:"Required; MaxSize(300)"`
	Course      *Course   `orm:"rel(fk); notnull"`
	Description string    `orm:"type(text); null"`
	Video       string    `orm:"size(300); null"`
	Metadata    string    `orm:"type(jsonb); null"`
	Views       int64     `orm:"default(0)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"auto_now;type(datetime)"`
}

type VideoLike struct {
	Id        int64         `orm:"auto"`
	Video     *CourseVideos `orm:"rel(fk); notnull"`
	User      *User         `orm:"rel(fk); notnull"`
	CreatedAt time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time     `orm:"auto_now;type(datetime)"`
}

type VideoComment struct {
	Id        int64         `orm:"auto"`
	Video     *CourseVideos `orm:"rel(fk); notnull"`
	User      *User         `orm:"rel(fk); notnull"`
	Comment   string        `orm:"size(500); notnull" valid:"Required; MaxSize(500)"`
	CreatedAt time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time     `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseVideos))
	orm.RegisterModel(new(VideoLike))
	orm.RegisterModel(new(VideoComment))
}

func GetAllVideoComments(videoID int64, page int) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var comments []*VideoComment

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("video_comment").Filter("video_id", videoID),
	}
	_, errs := query.Paginate().RelatedSel("user").OrderBy("-Id").All(&comments)

	if errs != nil {
		return nil, errs
	}

	data, err := query.CreatePagination(comments)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type CommentVideoSerializer struct {
	VideoID int64  `json:"video_id" valid:"Required"`
	Comment string `json:"comment" valid:"Required"`
}

func CommentVideo(videoID int64, user *User, comment string) (*VideoComment, error) {
	o := orm.NewOrm()

	commentObj := &VideoComment{
		Video:   &CourseVideos{Id: videoID},
		User:    user,
		Comment: comment,
	}

	if _, err := o.Insert(commentObj); err != nil {
		return nil, err
	}

	return commentObj, nil
}

type LikeVideoSerializer struct {
	VideoID int64 `json:"video_id" valid:"Required"`
}

func LikeVideo(videoID int64, user *User) (bool, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("video_like").Filter("video_id", videoID).Filter("user_id", user.Id)

	exists := qs.Exist()

	if exists {
		// Like exists → remove it (unlike)
		_, err := qs.Delete()
		if err != nil {
			return false, err
		}
		return false, nil
	}

	like := &VideoLike{
		Video: &CourseVideos{Id: videoID},
		User:  &User{Id: user.Id},
	}

	if _, err := o.Insert(like); err != nil {
		return false, err
	}

	return true, nil
}

func IncreaseViewCount(videoID int64) (int64, error) {
	o := orm.NewOrm()
	video := &CourseVideos{Id: videoID}

	if err := o.Read(video); err != nil {
		return 0, err
	}

	video.Views++
	if _, err := o.Update(video); err != nil {
		return 0, err
	}

	return video.Views, nil
}

type VideoLikeCountSerializer struct {
	UserLiked bool  `json:"user_liked"`
	TotalLike int64 `json:"total_like"`
}

func GetVideoLikeCount(videoID int64, user *User) (VideoLikeCountSerializer, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("video_like").Filter("video_id", videoID)
	count, err := qs.Count()
	if err != nil {
		return VideoLikeCountSerializer{}, err
	}
	userLiked := qs.Filter("user_id", user.Id).Exist()

	return VideoLikeCountSerializer{
		UserLiked: userLiked,
		TotalLike: count,
	}, nil
}
