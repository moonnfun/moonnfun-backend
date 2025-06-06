package global

import (
	"fmt"
	"log/slog"
)

func LogInit() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	switch Config.Log {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Debug(args ...any) {
	if Config.Log != "debug" {
		return
	}
	fmt.Printf("======debug-start======\n")
	strLog := ""
	for i, data := range args {
		if i == len(args)-1 {
			strLog += fmt.Sprintf("%+v", data)
		} else {
			strLog += fmt.Sprintf("%+v ", data)
		}
	}
	fmt.Println(strLog)
	fmt.Printf("======debug-end======\n")
}

func DebugForce(args ...any) {
	fmt.Printf("======debug-start======\n")
	strLog := ""
	for i, data := range args {
		if i == len(args)-1 {
			strLog += fmt.Sprintf("%+v", data)
		} else {
			strLog += fmt.Sprintf("%+v ", data)
		}
	}
	fmt.Println(strLog)
	fmt.Printf("======debug-end======\n")
}

// func Log(f func()) {
// 	if Config.Mode == "dev" {
// 		if f != nil {
// 			f()
// 		}
// 	}
// }
