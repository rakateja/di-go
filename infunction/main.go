package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

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

func storeNewEntry(sqlDB *sql.DB, entity Following) error {
	stmt, err := sqlDB.Prepare("INSERT INTO following (username, full_name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(entity.Username, entity.FullName)
	if err != nil {
		return err
	}
	return nil
}

type Following struct {
	ID       int
	Username string
	FullName string
}

func NewFollowing(username, fullName string) Following {
	return Following{0, username, fullName}
}

func main() {
	sqlDB, err := newSql("localhost", "raka", "raka_pass")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	newFollowing := NewFollowing("root", "I'm Root!")
	if err := storeNewEntry(sqlDB, newFollowing); err != nil {
		log.Fatalf("Error %v", err)
	}
	log.Println("Done!")
}
