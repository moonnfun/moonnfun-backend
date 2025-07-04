package logic

import (
	"fmt"
	"meme3/global"
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

	KlineType_Init   = "kline_init"
	KlineType_Update = "kline_update"
)

var (
	klinesInit   sync.Map
	klinesUpdate sync.Map

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
		// "second_1",
		"minute_1",
		"minute_5",
		"minute_15",
		"hour_4",
		"day_1",
		// "week_1",
	}
}

func getTimePeroid(timePeroid string) time.Duration {
	switch timePeroid {
	case "second_1":
		return Second
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

func InitKLine(klineType string, tokenAddress, timePeroid string) {
	if klineType == KlineType_Init {
		klinesInit.Store(klinesKey(tokenAddress, timePeroid), nil)
	} else {
		klinesUpdate.Store(klinesKey(tokenAddress, timePeroid), nil)
	}
}

func GetKLine(klineType string, tokenAddress, timePeroid string) any {
	if klineType == KlineType_Update {
		v, ok := klinesUpdate.Load(klinesKey(tokenAddress, timePeroid))
		if !ok || v == nil {
			vi, _ := klinesInit.Load(klinesKey(tokenAddress, timePeroid))
			return vi
		}
		return v
	} else {
		v, ok := klinesInit.Load(klinesKey(tokenAddress, timePeroid))
		if !ok || v == nil {
			return nil
		}
		return v
	}
}

func SaveKLine(klineType string, tokenAddress, timePeroid string, kLine model.TokenOHLCV) {
	if klineType == KlineType_Update {
		klinesUpdate.Store(klinesKey(tokenAddress, timePeroid), kLine)
	} else {
		klinesInit.Store(klinesKey(tokenAddress, timePeroid), kLine)
	}
}

func UpdateKLine(klineType string, tokenAddress string, timePeroid string, price decimal.Decimal, timestamp time.Time, volume decimal.Decimal) (bool, model.TokenOHLCV, error) {
	kline := GetKLine(klineType, tokenAddress, timePeroid)
	timeNow := getTimestamp(timePeroid, timestamp)

	bNext := false
	if kline != nil {
		startAt := getTimestamp(timePeroid, kline.(model.TokenOHLCV).StartAt)
		bNext = startAt.Equal(timeNow) == false
		if startAt.UnixMilli() > timeNow.UnixMilli() {
			return false, model.TokenOHLCV{}, fmt.Errorf("ingore kline, price: %v, timestamp: %v", price, timestamp.UnixMilli())
		}
		if timePeroid == "minute_1" {
			global.Debug("compare time: ", kline.(model.TokenOHLCV).StartAt.Unix(), timestamp.Unix(), startAt.Unix(), timeNow.Unix(), bNext)
		}
	}
	if kline == nil || bNext {
		// // 更新收盘价
		// if currentKLine != nil {
		// 	currentKLine.C = price
		// }

		// 启动新K线
		currentKLine := model.TokenOHLCV{
			T:       uint64(timeNow.UnixMilli()),
			OT:      uint64(timestamp.UnixMilli()),
			StartAt: timeNow,
		}
		if kline == nil {
			currentKLine.O = price
			currentKLine.H = currentKLine.O
			currentKLine.L = currentKLine.O
			currentKLine.C = currentKLine.O
			currentKLine.V = volume
			global.Debug("create: ", currentKLine)
		} else {
			currentKLine.O = kline.(model.TokenOHLCV).C
			currentKLine.H = currentKLine.O
			currentKLine.L = currentKLine.O
			currentKLine.C = price
			currentKLine.V = volume
			if price.GreaterThan(currentKLine.H) {
				currentKLine.H = price
			}
			if price.LessThan(currentKLine.L) {
				currentKLine.L = price
			}
			global.Debug("create next: ", kline, currentKLine)
		}
		SaveKLine(klineType, tokenAddress, timePeroid, currentKLine)
		return false, currentKLine, nil
	} else {
		currentKLine := kline.(model.TokenOHLCV)
		global.Debug("update1: ", kline)
		// 更新当前K线
		if price.GreaterThan(currentKLine.H) {
			currentKLine.H = price
		}
		if price.LessThan(currentKLine.L) {
			currentKLine.L = price
		}
		currentKLine.C = price
		// currentKLine.T = uint64(timestamp.UnixMilli())
		currentKLine.V = currentKLine.V.Add(volume)
		global.Debug("update2: ", currentKLine)
		SaveKLine(klineType, tokenAddress, timePeroid, currentKLine)
		return true, currentKLine, nil
	}
}

func getTimestamp(timePeroid string, t time.Time) time.Time {
	switch timePeroid {
	case "second_1":
		return t
	case "minute_1":
		return t.Truncate(time.Minute)
	case "minute_5":
		return t.Truncate(5 * time.Minute)
	case "minute_15":
		return t.Truncate(15 * time.Minute)
	case "hour_4":
		return t.Truncate(4 * time.Hour)
	case "day_1":
		y, m, d := t.UTC().Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	case "week_1":
		utc := t.UTC()
		weekday := int(utc.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7 instead of 0
		}
		startOfWeek := utc.AddDate(0, 0, -weekday+1)
		y, m, d := startOfWeek.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	default:
		return t
	}
}
