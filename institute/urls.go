package institute

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"] = append(web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"],
		web.ControllerComments{
			Method:           "Home",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"] = append(web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"],
		web.ControllerComments{
			Method:           "GetAllHomeBanner",
			Router:           "/banner",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"] = append(web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"],
		web.ControllerComments{
			Method:           "Get",
			Router:           "/:uid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"] = append(web.GlobalControllerRouter["allcoaching-go/institute:InstituteController"],
		web.ControllerComments{
			Method:           "GetAllCategories",
			Router:           "/categories",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
