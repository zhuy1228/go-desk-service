package services

import (
	"go-desk-service/libs"
	"go-desk-service/models"
	"time"

	"gorm.io/gorm"
)

type ProfileService struct {
	db *gorm.DB
}

func InitProfileService() *ProfileService {
	db := libs.GetDB()
	return &ProfileService{
		db: db,
	}
}

// 获取个人信息
func (p *ProfileService) GetInfo(params *models.Profile) models.Profile {
	var user models.Profile
	p.db.Model(&models.Profile{}).Where(params).First(&user)
	return user
}

// 更新账号信息
func (p *ProfileService) UpdateAccount(userId int64, params *models.Profile) bool {
	params.UpdatedAt = time.Now()
	p.db.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(params)
	return true
}

func (p *ProfileService) Create(params *models.Profile) bool {
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	p.db.Create(params)
	return true
}
