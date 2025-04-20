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

}
