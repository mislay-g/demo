package engine

import (
	"demo/src/memcache"
	"demo/src/models"
	"github.com/toolkits/pkg/container/list"
)

type WorkersType struct {
	pairs map[string]*PairsEval
}

var Workers = &WorkersType{
	pairs: make(map[string]*PairsEval),
}

func (ws *WorkersType) Build(pairIds []uint) {
	newPairs := make(map[string]*models.TradingPairs)

	for i := 0; i < len(pairIds); i++ {
		pair := memcache.TradingPairsCache.Get(pairIds[i])
		if pair == nil {
			continue
		}
		hash := pair.Symbol
		newPairs[hash] = pair
	}

	//stop old
	for hash := range Workers.pairs {
		if _, has := newPairs[hash]; !has {
			Workers.pairs[hash].Stop()
			delete(Workers.pairs, hash)
		}
	}

	//start new
	for hash := range newPairs {
		if _, has := Workers.pairs[hash]; has {
			//already exist
			continue
		}
		pe := &PairsEval{
			Channel:     newPairs[hash].Channel,
			Symbol:      newPairs[hash].Symbol,
			listenQuit:  make(chan struct{}),
			consumeQuit: make(chan struct{}),
			TradeQueue:  list.NewSafeListLimited(100000),
			MinQueue:    list.NewSafeListLimited(100000),
		}
		go pe.Start()
		Workers.pairs[hash] = pe
	}
}
