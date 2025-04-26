package users

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/users:UserController"] = append(web.GlobalControllerRouter["allcoaching-go/users:UserController"],
		web.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/users:UserController"] = append(web.GlobalControllerRouter["allcoaching-go/users:UserController"],
		web.ControllerComments{
			Method:           "Get",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/users:UserController"] = append(web.GlobalControllerRouter["allcoaching-go/users:UserController"],
		web.ControllerComments{
			Method:           "LoginUser",
			Router:           "/login",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
