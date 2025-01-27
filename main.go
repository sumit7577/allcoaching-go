package main

import (
	"allcoaching-go/allCoachingProject"
	_ "allcoaching-go/allCoachingProject"
	_ "allcoaching-go/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	orm.Debug = true
	allCoachingProject.SetDatabase()
	orm.RunSyncdb("default", false, true)
	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.Run()
}
