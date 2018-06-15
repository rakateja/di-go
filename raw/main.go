package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlHost     = "localhost"
	mysqlUser     = "root"
	mysqlPassword = "root-pass"
)

type Following struct {
	ID       int
	Username string
	FullName string
}

func main() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:3306)/following?parseTime=true", mysqlUser, mysqlPassword, mysqlHost)
	sqlDB, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error %v", err)
	}
	stmt, err := sqlDB.Prepare("INSERT INTO following (username, full_name) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	defer stmt.Close()

	following := NewFollowing("root", "I'm Root")
	_, err = stmt.Exec(following.Username, following.FullName)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
}
