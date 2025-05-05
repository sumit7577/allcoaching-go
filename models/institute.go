package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Institute represents the institute model.
type Institute struct {
	Id           int64     `orm:"auto"`
	Name         string    `orm:"size(100); null" valid:"MaxSize(100)"`
	About        string    `orm:"type(text); null"`
	Category     *Category `orm:"rel(fk); null"`
	Banner       *Banner   `orm:"rel(fk); null"`
	DirectorName string    `orm:"size(150); notnull" valid:"Required; MaxSize(150)"`
	User         *User     `orm:"rel(one); unique; notnull"`
	Followers    []*User   `orm:"rel(m2m); null"`
	Image        string    `orm:"size(300); null"`
	DateCreated  time.Time `orm:"auto_now_add;type(datetime)"`
	DateUpdated  time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Institute))
	//InitSearchVector()
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
	Course     []*Course       `json:"course"`
	Videos     []*CourseVideos `json:"videos"`
	TestSeries []*TestSeries   `json:"testseries"`
	Institute  *Institute      `json:"institute"`
	Documents  []*Documents    `json:"documents"`
}

func GetInstitute(uid int64, page int) (*PaginationSerializer, error) {
	o := orm.NewOrm()
	var courses []*Course
	var videos []*CourseVideos
	institute := Institute{}
	var testSeries []*TestSeries
	var documents []*Documents

	query := &Pagination{
		Offset: page,
		Limit:  10,
		query:  o.QueryTable("course").Filter("Institute__Id", uid),
	}

	//Fetch Course
	_, err := query.Paginate().RelatedSel("category").All(&courses)

	//Fetch one Institute all course videos
	_, err = o.QueryTable("course_videos").Filter("Course__Institute__Id", uid).Offset(page).Limit(10).All(&videos)

	//Fetch one Institute all course test series
	_, err = o.QueryTable("test_series").Filter("Course__Institute__Id", uid).Offset(page).Limit(10).All(&testSeries, "id", "name", "description", "questions", "timer", "created_at", "updated_at")

	_, err = o.QueryTable("documents").Filter("Course__Institute__Id", uid).Offset(page).Limit(10).All(&documents, "id", "name", "description", "file", "created_at", "updated_at")

	//Fetch institute detail
	err = o.QueryTable("institute").RelatedSel("banner", "category").Filter("Id", uid).One(&institute)

	if err != nil {
		return nil, err
	}

	serializer := &CourseSerializer{
		Course:     courses,
		Videos:     videos,
		Institute:  &institute,
		TestSeries: testSeries,
		Documents:  documents,
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
		_, err := o.QueryTable("institute").Filter("Category__Id", category.Id).Exclude("name", "Allcoaching").All(&institutes, "id", "name", "about", "director_name", "image", "date_created", "date_updated")
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

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	var categories []*Category
	_, err := o.QueryTable("category").All(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetAllHomeBanner() ([]*Banner, error) {
	o := orm.NewOrm()
	var course Course
	err := o.QueryTable("course").
		Filter("Institute__Name", "Allcoaching").
		Filter("Name", "Allcoaching_banner").
		One(&course)

	if course.Id == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	var banners []*Banner
	_, err = o.LoadRelated(&course, "Banner")

	if err != nil {
		return nil, nil
	}
	banners = append(banners, course.Banner...)

	return banners, nil
}

func GetSearchIns(query string) ([]*Institute, error) {
	o := orm.NewOrm()
	//InitSearchVector()
	var institutes []*Institute

	// Prepare the tsquery string with prefix search
	terms := strings.Fields(query)
	for i, term := range terms {
		terms[i] = term + ":*"
	}
	tsQuery := strings.Join(terms, " & ")

	sql := `
		SELECT id, name, about, image, director_name, date_created, date_updated
		FROM institute, to_tsquery('english', ?) query
		WHERE search_vector @@ query
		ORDER BY ts_rank(search_vector, query) DESC;
	`

	_, err := o.Raw(sql, tsQuery).QueryRows(&institutes)
	if err != nil {
		return nil, err
	}
	return institutes, nil
}

func InitSearchVector() {
	o := orm.NewOrm()

	sqls := []string{
		// 1. Add search_vector column if it doesn't exist
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name='institute' AND column_name='search_vector'
			) THEN
				ALTER TABLE institute ADD COLUMN search_vector tsvector;
			END IF;
		END
		$$;`,

		// 2. Populate search_vector initially
		`UPDATE institute
		 SET search_vector = to_tsvector('english',
			coalesce(name, '') || ' ' || coalesce(about, '') || ' ' || coalesce(director_name, ''));`,

		// 3. Create GIN index if not exists
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes
				WHERE indexname = 'idx_institute_search_vector'
			) THEN
				CREATE INDEX idx_institute_search_vector ON institute USING GIN(search_vector);
			END IF;
		END
		$$;`,

		// 4. Create trigger function if not exists
		`CREATE OR REPLACE FUNCTION institute_search_vector_trigger() RETURNS trigger AS $$
		begin
			new.search_vector :=
				to_tsvector('english',
				coalesce(new.name, '') || ' ' || coalesce(new.about, '') || ' ' || coalesce(new.director_name, ''));
			return new;
		end
		$$ LANGUAGE plpgsql;`,

		// 5. Create trigger itself
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_trigger WHERE tgname = 'tsvectorupdate'
			) THEN
				CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
				ON institute FOR EACH ROW EXECUTE FUNCTION institute_search_vector_trigger();
			END IF;
		END
		$$;`,
	}

	for _, sql := range sqls {
		_, err := o.Raw(sql).Exec()
		if err != nil {
			log.Println("❌ Failed SQL:", sql)
			log.Println("Error:", err)
		} else {
			log.Println("✅ Successfully executed SQL step.")
		}
	}
}

func ToggleFollowInstitute(uid int64, user *User) (string, error) {
	o := orm.NewOrm()
	institute := &Institute{Id: uid}

	count, err := o.QueryTable("institute_users").Filter("institute_id", uid).Filter("user_id", user.Id).Count()
	if err != nil {
		return fmt.Sprintf("Error checking follow status: %s", err), err
	}

	m2m := o.QueryM2M(institute, "Followers")
	if count > 0 {
		if _, err := m2m.Remove(user); err != nil {
			return "", err
		}
		return "Institute Unfollowed", nil
	} else {
		if _, err := m2m.Add(user); err != nil {
			return "", err
		}
		return "Institute Followed", nil
	}
}
