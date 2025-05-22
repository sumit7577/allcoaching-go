package main

import (
	"allcoaching-go/allcoachingProject"
	_ "allcoaching-go/allcoachingProject"
	_ "allcoaching-go/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	orm.Debug = true
	allcoachingProject.SetDatabase()
	orm.RunSyncdb("default", false, true)
	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.Run()
}
