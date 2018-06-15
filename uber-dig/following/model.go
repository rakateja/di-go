package following

import (
	"database/sql"

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
	Store(entity Following) error
}

type MysqlRepository struct {
	sqlDB *sql.DB
}

func NewMysqlRepository(sqlDB *sql.DB) Repository {
	return MysqlRepository{sqlDB: sqlDB}
}

func (sql MysqlRepository) Store(entity Following) error {
	stmt, err := sql.sqlDB.Prepare("INSERT INTO following (username, full_name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(entity.Username, entity.FullName); err != nil {
		return err
	}
	return nil
}

type Service struct {
	followingRepo Repository
}

func NewService(followingRepo Repository) Service {
	return Service{followingRepo}
}

func (s Service) Create(username, fullName string) error {
	err := s.followingRepo.Store(NewFollowing(username, fullName))
	if err != nil {
		return err
	}
	return nil
}
