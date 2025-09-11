package services

import (
	"go-desk-service/libs"
	"go-desk-service/models"

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
func (p *ProfileService) UpdateAccount(userId int8, params *models.Profile) bool {
	p.db.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(params)
	return true
}
