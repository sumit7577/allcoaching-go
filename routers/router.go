package routers

import (
	"allcoaching-go/institute"
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
	)
	web.AddNamespace(ns)
}
