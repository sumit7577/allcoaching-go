package allCoachingProject

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

func SetDatabase() {
	/*orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./data/allcoaching.db")*/
	orm.RegisterDriver("postgres", orm.DRPostgres)
	value, err := web.AppConfig.String("database-prod::dsn")
	if err != nil {
		panic(err)
	}
	orm.RegisterDataBase("default", "postgres", value)
}
