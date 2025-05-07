package videos

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/videos:VideoController"] = append(web.GlobalControllerRouter["allcoaching-go/videos:VideoController"],
		web.ControllerComments{
			Method:           "GetAllVideoComments",
			Router:           "/comments/:uid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/videos:VideoController"] = append(web.GlobalControllerRouter["allcoaching-go/videos:VideoController"],
		web.ControllerComments{
			Method:           "LikeVideo",
			Router:           "/like/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/videos:VideoController"] = append(web.GlobalControllerRouter["allcoaching-go/videos:VideoController"],
		web.ControllerComments{
			Method:           "CreateVideoComment",
			Router:           "/comment/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}
