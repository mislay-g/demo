package models

import (
	"github.com/pkg/errors"
)

type KLines struct {
	Id      uint    `json:"id" gorm:"primaryKey"`
	Symbol  string  `json:"symbol"`
	Channel string  `json:"channel"`
	Period  string  `json:"period"`
	Time    int64   `json:"time"`
	Open    float64 `json:"open"`
	Close   float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Num     float64 `json:"num"`
}

func (k *KLines) TableName() string {
	return "k_lines"
}

func (k *KLines) Add() error {
	return Insert(k)
}

func (k *KLines) Update(selectField interface{}, selectFields ...interface{}) error {
	return DB().Model(k).Select(selectField, selectFields...).Updates(k).Error
}

func (k *KLines) Del() error {
	return DB().Where("id=?", k.Id).Delete(&KLines{}).Error
}

func KLinesGets(symbol, period, begin string, page, size int) ([]*KLines, error) {
	session := DB().Limit(size).Offset(page * size)

	var ks []*KLines
	session.Where("symbol = ? and period = ? and `time` >= ?", symbol, period, begin)
	err := session.Find(&ks).Error
	if err != nil {
		return ks, errors.WithMessage(err, "failed to query kLines")
	}

	return ks, nil
}
