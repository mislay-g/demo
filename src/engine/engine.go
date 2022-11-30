package engine

import (
	"context"
	"demo/src/logger"
	"demo/src/memcache"
	"demo/src/models"
	"demo/src/storage"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/toolkits/pkg/container/list"
	"math"
	"time"
)

func Start() {
	for {
		duration := 9 * time.Second
		select {
		case <-context.Background().Done():
			return
		case <-time.After(duration):
			ids := memcache.TradingPairsCache.GetPairIds()
			Workers.Build(ids)
		}
	}
}

var C = cache.New(30*time.Minute, 30*time.Minute)

type PairsEval struct {
	Channel     string
	Symbol      string
	TradeQueue  *list.SafeListLimited
	MinQueue    *list.SafeListLimited
	listenQuit  chan struct{}
	consumeQuit chan struct{}
	publishQuit chan struct{}
}

func (p *PairsEval) Stop() {
	logger.Info("PairsEval:%d stopping" + p.Channel)
	go close(p.listenQuit)
	go close(p.consumeQuit)
	go close(p.publishQuit)
}

func (p *PairsEval) Start() {
	go listenTrade(p)
	go consumeTrade(p)
	go publishKLine(p)
}
func publishKLine(p *PairsEval) {
	duration := 500 * time.Millisecond
	for {
		select {
		case <-p.publishQuit:
			return
		case <-time.After(duration):
			if msg, exist := C.Get(fmt.Sprintf("KLines_%s_latest", p.Symbol)); exist {
				marshal, _ := json.Marshal(msg)
				storage.Redis.Publish(context.Background(), fmt.Sprintf("market.klines.%s.lasted", p.Symbol), marshal)
			}
		}
	}
}

func consumeTrade(p *PairsEval) {
	duration := time.Minute
	for {
		select {
		case <-p.consumeQuit:
			return
		case <-time.After(duration):
			consume(p)
		}
	}
}

func consume(p *PairsEval) {
	t := time.Now().Unix()
	l := p.TradeQueue.Len()
	bar := &Bar{
		Time:    t,
		Symbol:  p.Symbol,
		Channel: p.Channel,
		Open:    0,
		Close:   0,
		High:    0,
		Low:     0,
		//用0表示无数据
		Num: 0,
	}
	for i := 0; i < l; i++ {
		tick := p.TradeQueue.PopBack().(*Tick)
		if i == 0 {
			bar.Open = tick.Price
			bar.High = tick.Price
			bar.Low = tick.Price
		} else {
			bar.High = math.Max(bar.High, tick.Price)
			bar.Low = math.Min(bar.High, tick.Price)

		}
		bar.Close = tick.Price
		bar.Num += tick.Num
	}
	saveToDataBase(bar, "1min")
	p.MinQueue.PushFront(bar)
	if p.MinQueue.Len() == 5 {
		//这里是单线程，不用考虑安全
		bar = &Bar{
			Time:    t,
			Symbol:  p.Symbol,
			Channel: p.Channel,
			Open:    0,
			Close:   0,
			High:    0,
			Low:     0,
			//用0表示无数据
			Num: 0,
		}
		for i := 0; i < 5; i++ {
			tick := p.MinQueue.PopBack().(*Bar)
			if i == 0 {
				bar.Open = tick.Open
				bar.High = tick.High
				bar.Low = tick.Low
			} else {
				bar.High = math.Max(bar.High, tick.High)
				bar.Low = math.Min(bar.High, tick.Low)
			}
			bar.Close = tick.Close
			bar.Num += tick.Num
		}
		saveToDataBase(bar, "5min")
	}
}

func saveToDataBase(b *Bar, period string) {
	k := &models.KLines{
		Symbol:  b.Symbol,
		Channel: b.Channel,
		Period:  period,
		Time:    b.Time,
		Open:    b.Open,
		Close:   b.Close,
		High:    b.High,
		Low:     b.Low,
		Num:     b.Num,
	}
	_ = k.Add()
	C.Set(fmt.Sprintf("KLines_%s_latest", b.Symbol), k, 30*time.Minute)
}

func listenTrade(p *PairsEval) {
	pubsub := storage.Redis.Subscribe(context.Background(), p.Channel)
	defer pubsub.Close()

	//监视停止通道
	go func() {
		<-p.listenQuit
		pubsub.Close()
	}()

	ch := pubsub.Channel()
	//会随着pubsub自动关闭
	for msg := range ch {
		var t Tick
		err := json.Unmarshal([]byte(msg.Payload), &t)
		if err != nil {
			logger.Error("tick unmarshal failed", err)
		}
		p.TradeQueue.PushFront(&t)
	}
	fmt.Println(p.Channel + "end")
}

type Tick struct {
	Symbol    string
	Type      int
	Price     float64
	Num       float64
	Timestamp int64
}

type Bar struct {
	Time    int64
	Channel string
	Symbol  string
	Open    float64
	Close   float64
	High    float64
	Low     float64
	Num     float64
}
