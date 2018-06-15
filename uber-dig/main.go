package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	config "github.com/rakateja/di-go/uber-dig/config"
	"github.com/rakateja/di-go/uber-dig/following"
	mysql "github.com/rakateja/di-go/uber-dig/mysql"
	"go.uber.org/dig"
)

func buildContainer() *dig.Container {
	c := dig.New()
	c.Provide(func() config.Config {
		var cfg config.Config
		flag.StringVar(&cfg.Prefix, "prefix", "[following_service]", "log prefix")
		flag.IntVar(&cfg.Port, "port", 8080, "service's port")
		flag.StringVar(&cfg.MysqlHost, "mysql_host", "localhost", "mysql's host")
		flag.StringVar(&cfg.MysqlUser, "mysql_user", "raka", "mysql's user")
		flag.StringVar(&cfg.MysqlPassword, "mysql_password", "raka_pass", "mysql user's password")
		flag.Parse()
		return cfg
	})
	c.Provide(func(cfg config.Config) *log.Logger {
		return log.New(os.Stdout, cfg.Prefix, 0)
	})
	c.Provide(func(cfg config.Config, l *log.Logger) (*sql.DB, error) {
		return mysql.New(cfg)
	})
	c.Provide(func(sqlDB *sql.DB) following.Repository {
		return following.NewMysqlRepository(sqlDB)
	})
	c.Provide(func(repo following.Repository) following.Service {
		return following.NewService(repo)
	})
	return c
}

func server(l *log.Logger, sqlDB *sql.DB, service following.Service) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		l.Printf("Got request, URL %s", req.URL)
		err := sqlDB.Ping()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = service.Create("root", "I'm Root!")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Done!"))
	})
	l.Println("listening on port :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	c := buildContainer()
	if err := c.Invoke(server); err != nil {
		panic(err)
	}
}
