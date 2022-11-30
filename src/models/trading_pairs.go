package models

type TradingPairs struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Symbol      string `json:"symbol"`
	Channel     string `json:"channel"`
	Used        int    `json:"used"`
	Deleted     int    `json:"deleted"`
	CreatedTime int64  `json:"createdTime" gorm:"autoUpdateTime:milli"`
	UpdatedTime int64  `json:"updatedTime" gorm:"autoUpdateTime:milli"`
}

func (t *TradingPairs) TableName() string {
	return "trading_pairs"
}

func (t *TradingPairs) Add() error {
	return Insert(t)
}

func (t *TradingPairs) Update(selectField interface{}, selectFields ...interface{}) error {
	return DB().Model(t).Select(selectField, selectFields...).Updates(t).Error
}

func (t *TradingPairs) Del() error {
	return DB().Where("id=?", t.Id).Delete(&TradingPairs{}).Error
}

func TradingPairsGets(where string, args ...interface{}) ([]*TradingPairs, error) {
	var lst []*TradingPairs
	err := DB().Where(where, args...).Find(&lst).Error
	return lst, err
}

func TradingPairsStatistics() (*Statistics, error) {
	session := DB().Model(&TradingPairs{}).Select("count(*) as total", "max(updated_time) as last_updated").Where("used = ? and deleted = ?", 1, 0)

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

type Statistics struct {
	Total       int64 `gorm:"total"`
	LastUpdated int64 `gorm:"last_updated"`
}
