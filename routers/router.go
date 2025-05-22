package routers

import (
	"allcoaching-go/courses"
	"allcoaching-go/institute"
	"allcoaching-go/orders"
	"allcoaching-go/posts"
	"allcoaching-go/test_series"
	"allcoaching-go/users"
	"allcoaching-go/videos"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	ns := web.NewNamespace("/v1",

		web.NSNamespace("/user",
			web.NSInclude(
				&users.UserController{},
			),
		),
		web.NSNamespace("/institute",
			web.NSInclude(
				&institute.InstituteController{},
			),
		),
		web.NSNamespace("/course",
			web.NSInclude(
				&courses.CourseController{},
			),
		),
		web.NSNamespace("/testSeries",
			web.NSInclude(
				&test_series.TestSeriesController{},
			),
		),
		web.NSNamespace("/orders",
			web.NSInclude(
				&orders.OrderController{},
			),
		),
		web.NSNamespace("/posts",
			web.NSInclude(
				&posts.PostController{},
			),
		),
		web.NSNamespace("/videos",
			web.NSInclude(
				&videos.VideoController{},
			),
		),
	)
	web.AddNamespace(ns)
}
