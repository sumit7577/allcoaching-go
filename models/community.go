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

type CommunityLike struct {
	Id        int64          `orm:"auto"`
	Post      *CommunityPost `orm:"rel(fk); notnull"`
	User      *User          `orm:"rel(fk); notnull"`
	CreatedAt time.Time      `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time      `orm:"auto_now;type(datetime)"`
}

type CommunityComment struct {
	Id        int64          `orm:"auto"`
	Post      *CommunityPost `orm:"rel(fk); notnull"`
	User      *User          `orm:"rel(fk); notnull"`
	Comment   string         `orm:"size(500); notnull" valid:"Required; MaxSize(500)"`
	CreatedAt time.Time      `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time      `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CommunityPost))
	orm.RegisterModel(new(CommunityLike))
	orm.RegisterModel(new(CommunityComment))
}

type CommunityPostWithMeta struct {
	Post         *CommunityPost `json:"post"`     // Embedding the original CommunityPost struct
	LikeCount    int64          `json:"likes"`    // Field for the number of likes
	CommentCount int64          `json:"comments"` // Field for the number of comments
	UserLiked    bool           `json:"user_liked"`
}

type LikeResult struct {
	PostID int64 `orm:"column(post_id)"`
	Count  int64 `orm:"column(count)"`
}

type CommentResult struct {
	PostID int64 `orm:"column(post_id)"`
	Count  int64 `orm:"column(count)"`
}

func GetAllPosts(page int, user *User) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var posts []*CommunityPost

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("community_post"),
	}
	_, errs := query.Paginate().RelatedSel("institute").OrderBy("-Id").All(&posts)

	if errs != nil {
		return nil, errs
	}

	var postIDs []int64
	for _, post := range posts {
		postIDs = append(postIDs, post.Id)
	}

	// Fetch like counts grouped by post ID

	var likeResults []LikeResult

	_, err := o.QueryTable("community_like").
		Filter("post__in", postIDs).
		GroupBy("post_id").
		Aggregate("post_id, COUNT(*) as count").
		All(&likeResults)

	if err != nil {
		return nil, err
	}

	likeMap := make(map[int64]int64)
	for _, lr := range likeResults {
		likeMap[lr.PostID] = lr.Count
	}

	var commentResults []CommentResult
	_, err = o.QueryTable("community_comment").
		Filter("post__in", postIDs).
		GroupBy("post_id").
		Aggregate("post_id, COUNT(*) as count").
		All(&commentResults)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[int64]int64)
	for _, cr := range commentResults {
		commentMap[cr.PostID] = cr.Count
	}

	var likedPostIDs []orm.Params
	_, err = o.QueryTable("community_like").
		Filter("user_id", user.Id).
		Filter("post_id__in", postIDs).
		Values(&likedPostIDs, "post_id")

	likedMap := make(map[int64]bool)
	for _, lp := range likedPostIDs {
		if postID, ok := lp["Post__Post"].(int64); ok {
			likedMap[postID] = true
		}
	}

	if err != nil {
		return nil, err
	}

	var enrichedPosts []*CommunityPostWithMeta
	for _, post := range posts {
		enriched := &CommunityPostWithMeta{
			Post:         post,
			LikeCount:    likeMap[post.Id],
			CommentCount: commentMap[post.Id],
			UserLiked:    likeMap[post.Id] > 0,
		}
		enrichedPosts = append(enrichedPosts, enriched)
	}
	data, err := query.CreatePagination(enrichedPosts)
	if err != nil {
		return nil, err
	}

	return data, nil

}

type LikePostSerializer struct {
	PostID int64 `json:"post_id" valid:"Required"`
}

func LikePost(postID int64, user *User) (bool, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("community_like").Filter("post_id", postID).Filter("user_id", user.Id)

	exists := qs.Exist()

	if exists {
		// Like exists → remove it (unlike)
		_, err := qs.Delete()
		if err != nil {
			return false, err
		}
		return false, nil
	}

	like := &CommunityLike{
		Post: &CommunityPost{Id: postID},
		User: &User{Id: user.Id},
	}

	if _, err := o.Insert(like); err != nil {
		return false, err
	}

	return true, nil

}

type CommentPostSerializer struct {
	PostID  int64  `json:"post_id" valid:"Required"`
	Comment string `json:"comment" valid:"Required"`
}

func CommentPost(postID int64, user *User, comment string) (*CommunityComment, error) {
	o := orm.NewOrm()

	commentObj := &CommunityComment{
		Post:    &CommunityPost{Id: postID},
		User:    user,
		Comment: comment,
	}

	if _, err := o.Insert(commentObj); err != nil {
		return nil, err
	}

	return commentObj, nil
}

func GetAllPostComments(postID int64, page int) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var comments []*CommunityComment

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("community_comment").Filter("post_id", postID),
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
