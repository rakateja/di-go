package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/dig"

	config "github.com/rakateja/di-go/uber-dig/config"
	"github.com/rakateja/di-go/uber-dig/following"
	mysql "github.com/rakateja/di-go/uber-dig/mysql"
)

func buildContainer() *dig.Container {
	c := dig.New()

	c.Provide(func() config.Config {
		var cfg config.Config
		flag.StringVar(&cfg.Prefix, "prefix", "[following_service] ", "log prefix")
		flag.IntVar(&cfg.Port, "port", 8080, "service's port")
		flag.StringVar(&cfg.MysqlHost, "mysql_host", "localhost", "mysql's host")
		flag.StringVar(&cfg.MysqlUser, "mysql_user", "root", "mysql's user")
		flag.StringVar(&cfg.MysqlPassword, "mysql_password", "root-pass", "mysql user's password")
		flag.StringVar(&cfg.MysqlDB, "mysql_db", "following", "mysql database")

		flag.Parse()
		return cfg
	})
	c.Provide(func(cfg config.Config) *log.Logger {
		return log.New(os.Stdout, cfg.Prefix, 0)
	})
	c.Provide(mysql.New)
	c.Provide(following.NewMysqlRepository)
	c.Provide(following.NewService)
	c.Provide(NewServer)
	return c
}

type Server struct {
	logger  *log.Logger
	sqlDB   *sql.DB
	service following.Service
}

func NewServer(l *log.Logger, sqlDB *sql.DB, service following.Service) Server {
	return Server{l, sqlDB, service}
}

func (s Server) Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		s.logger.Printf("Got request, URL %s", req.URL)
		err := s.sqlDB.Ping()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = s.service.Create("root", "I'm Root!")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Done!"))
	})
	s.logger.Println("listening on port :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	c := buildContainer()
	if err := c.Invoke(func(s Server) {
		s.Serve()
	}); err != nil {
		panic(err)
	}
}
