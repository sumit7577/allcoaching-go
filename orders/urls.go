package orders

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/orders:OrderController"] = append(web.GlobalControllerRouter["allcoaching-go/orders:OrderController"],
		web.ControllerComments{
			Method:           "GetOrderEvents",
			Router:           "/events",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}
