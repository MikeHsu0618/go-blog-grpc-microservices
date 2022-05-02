package dbcontext

import (
	"fmt"

	"blog-grpc-microservices/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(dsn string) (*DB, error) {
	params := " sslmode=disable TimeZone=Asia/Taipei"
	dsn = dsn + params
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func NewUserDB(conf *config.Config) (*DB, error) {
	// dsn := fmt.Sprintf(
	// 	"%s:%s@tcp(%s%s)/%s",
	// 	conf.User.DB.User,
	// 	conf.User.DB.Password,
	// 	conf.User.DB.Host,
	// 	conf.User.DB.Port,
	// 	conf.User.DB.Name,
	// )

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v",
		conf.User.DB.Host,
		conf.User.DB.User,
		conf.User.DB.Password,
		conf.User.DB.Name,
		conf.User.DB.Port,
	)
	return NewDB(dsn)
}

func NewPostDB(conf *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s",
		conf.Post.DB.User,
		conf.Post.DB.Password,
		conf.Post.DB.Host,
		conf.Post.DB.Port,
		conf.Post.DB.Name,
	)
	return NewDB(dsn)
}

func NewCommentDB(conf *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s",
		conf.Comment.DB.User,
		conf.Comment.DB.Password,
		conf.Comment.DB.Host,
		conf.Comment.DB.Port,
		conf.Comment.DB.Name,
	)
	return NewDB(dsn)
}
