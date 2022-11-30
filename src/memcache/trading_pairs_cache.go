package memcache

import (
	"demo/src/logger"
	"demo/src/models"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sync"
	"time"
)

type TradingPairsCacheType struct {
	statTotal       int64
	statLastUpdated int64

	sync.RWMutex
	pairs map[uint]*models.TradingPairs
}

var TradingPairsCache = TradingPairsCacheType{
	statTotal:       -1,
	statLastUpdated: -1,
	pairs:           make(map[uint]*models.TradingPairs),
}

func (tpc *TradingPairsCacheType) Reset() {
	tpc.Lock()
	defer tpc.Unlock()

	tpc.statTotal = -1
	tpc.statLastUpdated = -1
	tpc.pairs = make(map[uint]*models.TradingPairs)
}

func (tpc *TradingPairsCacheType) StatChanged(total, lastUpdated int64) bool {
	if tpc.statTotal == total && tpc.statLastUpdated == lastUpdated {
		return false
	}

	return true
}

func (tpc *TradingPairsCacheType) Set(m map[uint]*models.TradingPairs, total, lastUpdated int64) {
	tpc.Lock()
	tpc.pairs = m
	tpc.Unlock()

	tpc.statTotal = total
	tpc.statLastUpdated = lastUpdated
}

func (tpc *TradingPairsCacheType) Get(pairId uint) *models.TradingPairs {
	tpc.RLock()
	defer tpc.RUnlock()
	return tpc.pairs[pairId]
}

func (tpc *TradingPairsCacheType) GetPairIds() []uint {
	tpc.RLock()
	defer tpc.RUnlock()

	count := len(tpc.pairs)
	list := make([]uint, 0, count)
	for pairId := range tpc.pairs {
		list = append(list, pairId)
	}

	return list
}

func SyncTradingPairs() {
	err := syncTradingPairs()
	if err != nil {
		fmt.Println("failed to sync alert rules:", err)
		exit(1)
	}

	go loopSyncAlertRules()
}

func loopSyncAlertRules() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := syncTradingPairs(); err != nil {
			logger.Error("failed to sync alert rules:", err)
		}
	}
}

func syncTradingPairs() error {
	stat, err := models.TradingPairsStatistics()
	if err != nil {
		return errors.WithMessage(err, "failed to exec TradingPairsStatistics")
	}
	if !TradingPairsCache.StatChanged(stat.Total, stat.LastUpdated) {
		return nil
	}

	lst, err := models.TradingPairsGets("used = ? and deleted = ?", 1, 0)
	if err != nil {
		return errors.WithMessage(err, "failed to exec TradingPairsGets")
	}

	m := make(map[uint]*models.TradingPairs)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Id] = lst[i]
	}

	TradingPairsCache.Set(m, stat.Total, stat.LastUpdated)
	return nil
}

func exit(code int) {
	os.Exit(code)
}
