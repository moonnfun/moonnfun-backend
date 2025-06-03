package server

import (
	"context"
	"embed"
	"meme3/global"
	"meme3/server/api"
	"meme3/server/web"
)

func Init(fatalCheck func(...any)) {
	fatalCheck(initServer())
}

func RunWeb(ctx context.Context, f *embed.FS) error {
	originList := []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:5137"}
	if global.Config.HostURL != "" {
		originList = append(originList, global.Config.HostURL)
	}
	router := web.Init(global.Config.WebPath, f, api.RouterInit)
	go web.Run(ctx, global.Config.Port, router, false, originList...)
	return nil
}

func initServer() error {
	// session
	if err := web.InitSession(); err != nil {
		return err
	}
	return nil
}
