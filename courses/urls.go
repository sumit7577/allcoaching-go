package courses

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/courses:CourseController"] = append(web.GlobalControllerRouter["allcoaching-go/courses:CourseController"],
		web.ControllerComments{
			Method:           "Purchase",
			Router:           "/purchase/:uid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}
