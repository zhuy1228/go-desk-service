package models

import "time"

type Profile struct {
	ID          int8      `gorm:"id;primary_key"`
	UserId      int8      `gorm:"user_id"`
	Nickname    string    `gorm:"nickname"`
	Email       string    `gorm:"email"`
	Phone       string    `gorm:"phone"`
	Status      int       `gorm:"status"`
	LoginStatus int       `gorm:"login_status"`
	CreatedAt   time.Time `gorm:"created_at;type:timestamptz"`
	UpdatedAt   time.Time `gorm:"updated_at;type:timestamptz"`
}

// 实现 TableName 方法指定表名
func (Profile) TableName() string {
	return "gd_profile"
}
