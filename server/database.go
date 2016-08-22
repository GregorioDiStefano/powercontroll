package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database interface {
	openDatabase()
	addDevice(mac string) (bool, error)
	//findDevice(mac string) (bool, error)
	//editDevice(mac string, fields map[string]string) error
	updateDeviceLocation(mac, location string) error
	deleteDevice(mac string) error
	setDeviceStatus(mac, status string) error
	getDevices() []device
	closeDatabase()
}

type ormDatabase struct {
	dbHandle *gorm.DB
}

const (
	deviceNotFound     = "Device not found"
	deviceExists       = "This device exists"
	deviceDeleteFailed = "Failed to delete device"
	deviceUpdateFailed = "Failed to update device"
)

func (db *ormDatabase) doesDeviceExist(mac string) bool {
	var count int
	db.dbHandle.Model(&device{}).Where("device_mac = ?", mac).Count(&count)

	if count == 0 {
		return false
	}
	return true

}

func (db *ormDatabase) openDatabase() {
	orm, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	orm.AutoMigrate(&device{})
	db.dbHandle = orm
	fmt.Println(db.dbHandle, orm)
}

func (db ormDatabase) addDevice(mac string) (bool, error) {
	if db.doesDeviceExist(mac) {
		return false, nil
	}
	if err := db.dbHandle.Create(&device{DeviceMac: mac, Location: "", Date: time.Now()}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (db ormDatabase) getDevices() []device {
	var devices []device
	db.dbHandle.Find(&devices)
	return devices
}

func (db ormDatabase) closeDatabase() {
	db.dbHandle.Close()
}

func (db ormDatabase) updateDeviceLocation(mac, location string) error {
	if !db.doesDeviceExist(mac) {
		return errors.New(deviceNotFound)
	}

	if err := db.dbHandle.Model(&device{}).Where("device_mac = ?", mac).Update("location", location).Error; err != nil {
		fmt.Println(err.Error())
		return errors.New(deviceUpdateFailed)
	}
	return nil
}

func (db ormDatabase) deleteDevice(mac string) error {
	if !db.doesDeviceExist(mac) {
		return errors.New(deviceNotFound)
	}

	if err := db.dbHandle.Model(&device{}).Where("device_mac = ?", mac).Delete(&device{}).Error; err != nil {
		return errors.New(deviceDeleteFailed)
	}
	return nil
}

func (db ormDatabase) setDeviceStatus(mac, status string) error {
	var statusBool bool

	if !db.doesDeviceExist(mac) {
		return errors.New(deviceNotFound)
	}

	if strings.ToLower(status) == "on" {
		statusBool = true
	} else if strings.ToLower(status) == "off" {
		statusBool = false
	} else {
		panic("invalid status")
	}

	if err := db.dbHandle.Model(&device{}).Where("device_mac = ?", mac).Update("powered", statusBool).Error; err != nil {
		fmt.Println(err)
		return errors.New(deviceDeleteFailed)
	}
	return nil
}
