package logic

import (
	"fmt"
	"meme3/service/model"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Second                    = 1000 * Millisecond
)

var (
	klines sync.Map

	v_Lock_KLine = new(sync.Mutex)

	// v_time_list = map[string]time.Duration{
	// 	"second_1":  Second,
	// 	"second_5":  5 * Second,
	// 	"second_8":  8 * Second,
	// 	"second_15": 15 * Second,
	// 	"minute_1":  60 * Second,
	// 	"minute_5":  60 * 5 * Second,
	// 	"minute_30": 60 * 30 * Second,
	// 	"hour_1":    60 * 60 * Second,
	// 	"hour_4":    4 * 60 * 60 * Second,
	// 	"day_1":     24 * 60 * 60 * Second,
	// 	"week_1":    7 * 24 * 60 * 60 * Second,
	// 	"month_1":   30 * 24 * 60 * 60 * Second,
	// 	"year_1":    365 * 24 * 60 * 60 * Second,
	// }
)

func getTimeList() []string {
	return []string{
		"minute_1",
		"minute_15",
		"hour_4",
		"minute_15",
		"day_1",
		"week_1",
	}
}

func getTimePeroid(timePeroid string) time.Duration {
	switch timePeroid {
	case "minute_1":
		return 60 * Second
	case "minute_15":
		return 60 * 15 * Second
	case "hour_4":
		return 4 * 60 * 60 * Second
	case "day_1":
		return 24 * 60 * 60 * Second
	case "week_1":
		return 7 * 24 * 60 * 60 * Second
	default:
		return Second
	}
}

func klinesKey(tokenAddress, timePeroid string) string {
	return fmt.Sprintf("%s-%v", tokenAddress, timePeroid)
}

func GetKLine(tokenAddress, timePeroid string) any {
	v, ok := klines.Load(klinesKey(tokenAddress, timePeroid))
	if !ok {
		return nil
	}
	return v
}

func SaveKLine(tokenAddress, timePeroid string, kLine model.TokenOHLCV) {
	klines.Store(klinesKey(tokenAddress, timePeroid), kLine)
}

func UpdateKLine(tokenAddress string, timePeroid time.Duration, price decimal.Decimal, timestamp time.Time, volume decimal.Decimal) (model.TokenOHLCV, error) {
	kline := GetKLine(tokenAddress, timePeroid.String())

	bNext := false
	if kline != nil {
		bNext = timestamp.After(kline.(model.TokenOHLCV).StartAt.Add(timePeroid))
	}
	if kline == nil || bNext {
		// // 更新收盘价
		// if currentKLine != nil {
		// 	currentKLine.C = price
		// }

		// 启动新K线
		currentKLine := model.TokenOHLCV{
			O:       price,
			C:       price,
			H:       price,
			L:       price,
			V:       volume,
			T:       uint64(time.Now().UnixMilli()),
			StartAt: timestamp,
		}
		SaveKLine(tokenAddress, timePeroid.String(), currentKLine)
		return currentKLine, nil
	} else {
		currentKLine := kline.(model.TokenOHLCV)
		// 更新当前K线
		if price.GreaterThan(currentKLine.H) {
			currentKLine.H = price
		}
		if price.LessThan(currentKLine.L) {
			currentKLine.L = price
		}
		currentKLine.C = price
		currentKLine.V = currentKLine.V.Add(volume)
		return currentKLine, nil
	}
}
