package models

import (
	"demo/src/storage"
	"gorm.io/gorm"
)


func DB() *gorm.DB {
	return storage.DB
}

func Count(tx *gorm.DB) (int64, error) {
	var cnt int64
	err := tx.Count(&cnt).Error
	return cnt, err
}

func Exists(tx *gorm.DB) (bool, error) {
	num, err := Count(tx)
	return num > 0, err
}

func Insert(obj interface{}) error {
	return DB().Create(obj).Error
}
