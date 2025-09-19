package services

import (
	"go-desk-service/config"
	"go-desk-service/libs"
	"go-desk-service/models"
	"go-desk-service/utils"
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

func (d *DeviceService) UpdateByToken(token string, device *models.Device) bool {
	device.UpdatedAt = time.Now()
	d.db.Model(&models.Device{}).Where("token = ?", token).Updates(device)
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
	err := d.db.Model(&models.Device{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 1, nil // 用户已登录
	}
	return 0, nil // 用户未登录
}

func (d *DeviceService) SetUserLoginStatus(userId int64, token, deviceUid string) {
	var device models.Device
	err := d.db.Model(&models.Device{}).Where("user_id = ? AND device_uid = ?", userId, deviceUid).First(&device).Error
	if err != nil {
		// 没查到就创建
		appConfig, _ := config.LoadConfig()
		Snowflake, _ := utils.NewSnowflake(appConfig.WorkerID, appConfig.DatacenterID)
		DeviceId := Snowflake.NextID()
		d.db.Create(&models.Device{
			ID:          DeviceId,
			UserId:      userId,
			Token:       token,
			DeviceUid:   deviceUid,
			LoginStatus: 1,
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
		})
		return
	}
	d.db.Model(&models.Device{}).Where("id = ?", device.ID).Updates(&models.Device{
		Token:       token,
		LoginStatus: 1,
		UpdatedAt:   time.Now(),
	})
}

func (d *DeviceService) Logout(token string) {
	d.db.Model(&models.Device{}).Where("token = ?", token).Updates(map[string]interface{}{
		"login_status": 0,
		"updated_at":   time.Now(),
	})
}

func (d *DeviceService) GetAll(userId int64) []models.Device {
	var device []models.Device
	d.db.Select("id", "hostname", "platform", "platform_version", "mac", "cpu", "mem", "disk", "device_uid", "login_status", "ip").Model(&models.Device{}).Where("user_id = ?", userId).Find(&device)

	return device
}

func (d *DeviceService) Login(token string) models.Device {
	d.db.Model(&models.Device{}).Where("token = ?", token).Updates(&models.Device{
		LoginStatus: 1,
		UpdatedAt:   time.Now(),
	})
	var device models.Device
	d.db.Model(&models.Device{}).Where("token = ?", token).First(&device)
	return device
}
