package posts

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/posts:PostController"] = append(web.GlobalControllerRouter["allcoaching-go/posts:PostController"],
		web.ControllerComments{
			Method:           "Get",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/posts:PostController"] = append(web.GlobalControllerRouter["allcoaching-go/posts:PostController"],
		web.ControllerComments{
			Method:           "GetByInsID",
			Router:           "/ins/:uid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/posts:PostController"] = append(web.GlobalControllerRouter["allcoaching-go/posts:PostController"],
		web.ControllerComments{
			Method:           "GetAllPostComments",
			Router:           "/comments/:uid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/posts:PostController"] = append(web.GlobalControllerRouter["allcoaching-go/posts:PostController"],
		web.ControllerComments{
			Method:           "LikePost",
			Router:           "/like/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/posts:PostController"] = append(web.GlobalControllerRouter["allcoaching-go/posts:PostController"],
		web.ControllerComments{
			Method:           "CreatePostComment",
			Router:           "/comment/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
