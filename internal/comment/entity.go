package comment

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint64         `json:"id"`
	UUID      string         `json:"uuid"`
	UserID    uint64         `json:"user_id"`
	PostID    uint64         `json:"post_id"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
