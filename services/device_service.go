package services

import (
	"go-desk-service/libs"
	"go-desk-service/models"
	"time"

	"gorm.io/gorm"
)

type DeviceService struct {
	db *gorm.DB
}

func InitDeviceService() *DeviceService {
	db := libs.GetDB()
	return &DeviceService{
		db: db,
	}
}

func (d *DeviceService) Create(device *models.Device) bool {
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	d.db.Create(device)
	return true
}

func (d *DeviceService) GetByUserIdAndDeviceUid(userId int64, deviceUid string) (*models.Device, error) {
	var device models.Device
	err := d.db.Model(&models.Device{}).Where("user_id = ? AND device_uid = ?", userId, deviceUid).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (d *DeviceService) Update(device *models.Device) bool {
	device.UpdatedAt = time.Now()
	d.db.Model(&models.Device{}).Where("id = ?", device.ID).Updates(device)
	return true
}

func (d *DeviceService) GetInfoByToken(token string) (*models.Device, error) {
	var device models.Device
	err := d.db.Model(&models.Device{}).Where("token = ?", token).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (d *DeviceService) GetUserLoginStatus(userId int64) (int, error) {
	var count int64
	err := d.db.Model(&models.Device{}).Where("user_id = ? AND login_status = ?", userId, 1).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 1, nil // 用户已登录
	}
	return 0, nil // 用户未登录
}
