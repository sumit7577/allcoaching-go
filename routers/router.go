package routers

import (
	"allcoaching-go/institute"
	"allcoaching-go/test_series"
	"allcoaching-go/users"

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
		web.NSNamespace("/testSeries",
			web.NSInclude(
				&test_series.TestSeriesController{},
			),
		),
	)
	web.AddNamespace(ns)
}
