package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Following struct {
	ID       int
	Username string
	FullName string
}

func NewFollowing(username, fullName string) Following {
	return Following{0, username, fullName}
}

type Repository interface {
	Store(following Following) error
}

func NewMysqlRepository(sqlDB *sql.DB) Repository {
	return MysqlRepository{sqlDB: sqlDB}
}

type MysqlRepository struct {
	sqlDB *sql.DB
}

func (sql MysqlRepository) Store(following Following) error {
	stmt, err := sql.sqlDB.Prepare("INSERT INTO following (username, full_name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(following.Username, following.FullName); err != nil {
		return err
	}
	return nil
}

func newSql(mysqlHost, mysqlUser, mysqlPassword string) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:3306)/following", mysqlUser, mysqlPassword, mysqlHost)
	sqlDB, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return sqlDB, nil
}

type Service struct {
	Repository
}

func NewService(followingRepository Repository) Service {
	return Service{followingRepository}
}

func (service Service) Create(username, fullName string) error {
	err := service.Store(NewFollowing(username, fullName))
	if err != nil {
		return err
	}
	return nil
}

var (
	mysqlHost     = "localhost"
	mysqlUser     = "root"
	mysqlPassword = "root-pass"
)

func main() {
	sqlDB, err := newSql(mysqlHost, mysqlUser, mysqlPassword)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	mysqlRepository := NewMysqlRepository(sqlDB)
	followingService := NewService(mysqlRepository)
	if err := followingService.Create("root", "I'm Root!"); err != nil {
		log.Fatalf("Error %v", err)
	}
	log.Println("Done!")
}
