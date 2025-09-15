package models

import "time"

type Device struct {
	ID              int64     `gorm:"id;primary_key"`
	UserId          int64     `gorm:"user_id"`
	DeviceUid       string    `gorm:"device_uid"`
	Token           string    `gorm:"token"`
	Hostname        string    `gorm:"hostname"`
	Platform        string    `gorm:"platform"`
	PlatformVersion string    `gorm:"platform_version"`
	Mac             string    `gorm:"mac"`
	Cpu             string    `gorm:"cpu"`
	Mem             string    `gorm:"mem"`
	Disk            string    `gorm:"disk"`
	LoginStatus     int       `gorm:"login_status"` // 1: 已登录, 0: 未登录
	CreatedAt       time.Time `gorm:"created_at;type:timestamptz"`
	UpdatedAt       time.Time `gorm:"updated_at;type:timestamptz"`
}

// 实现 TableName 方法指定表名
func (Device) TableName() string {
	return "gd_device"
}
