package test_series

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"] = append(web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"],
		web.ControllerComments{
			Method:           "SubmitAnswer",
			Router:           "/attempt",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"] = append(web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"],
		web.ControllerComments{
			Method:           "Result",
			Router:           "/result",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"] = append(web.GlobalControllerRouter["allcoaching-go/test_series:TestSeriesController"],
		web.ControllerComments{
			Method:           "Submit",
			Router:           "/submit",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
