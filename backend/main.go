package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"meme3/global"
	"meme3/server"
	"meme3/server/wss"
	"meme3/service/logic"
	"meme3/service/monitor"
	"meme3/service/store"

	"os"
	"os/signal"
	"syscall"
)

var (
	f *embed.FS

	port    = flag.Int("port", 5137, "--port=5137")
	debug   = flag.Bool("debug", false, "--debug=true")
	logl    = flag.String("log", "debug", "--log=info")
	webPath = flag.String("webPath", "dist", "--webPath=dist")
	hostUrl = flag.String("hostUrl", "http://localhost:5137", "--host=https://www.example.com")

	cmd                    = flag.String("cmd", "", "--cmd=SyncForUpdate")
	refreshSortListTimer   = flag.Int("refreshSortListTimer", 1, "--refreshSortListTimer=1")
	refreshTokenPriceTimer = flag.Int("refreshTokenPriceTimer", 1, "--refreshTokenPriceTimer=1")

	launch   = flag.Float64("launch", 10, "--launch=100000000")
	node     = flag.String("node", "node", "--node=/usr/bin/node")
	rpc      = flag.String("rpc", "", "--rpc=/path/to/rpc.json")
	testnet  = flag.Bool("testnet", false, "--testnet=false")
	sprice   = flag.Float64("sprice", 0.000001, "--startPrice=0.000000001")
	dbName   = flag.String("dbName", "meme3-local", "--dbName=meme3")
	hotList  = flag.String("hotList", "", "--hostList=/usr/local/hotList.json")
	dbUrl    = flag.String("dbUrl", "mongodb://localhost:37015", "--log=mongodb://localhost:37015")
	contract = flag.String("contract", "", "--dbName=xxx")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	// init global
	Init(ctx)

	slog.Info("system init successed", "log", global.Config.Log)

	InitLogic(ctx)
	InitCache()
	runServer(ctx, cancel)

	// wait signal
	var sigch = make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) //, syscall.SIGUSR1, syscall.SIGUSR2)
	<-sigch
	cancel()
}

func runServer(ctx context.Context, cancel func()) {
	go monitor.WatchContractEvents(monitor.HandleTokenCreatedTx, monitor.HandleTokenTradingTx, monitor.HandleTokenSwapTx)

	server.Init(FatalCheck)
	FatalCheck(server.RunWeb(ctx, f))
}

func Init(ctx context.Context) {
	global.Config.Cmd = *cmd
	global.Config.Port = *port
	global.Config.Debug = *debug
	global.Config.WebPath = *webPath
	global.Config.HostURL = *hostUrl
	global.Config.Log = *logl
	global.Config.LaunchMC = *launch
	global.Config.NodePath = *node
	global.Config.Testnet = *testnet
	global.Config.StartPrice = *sprice
	global.Config.ContractAddress = *contract
	global.Config.RefreshSortListTimer = *refreshSortListTimer
	global.Config.RefreshTokenPriceTimer = *refreshTokenPriceTimer
	global.Config.WebSessionTimeout = global.C_websession_timeout
	global.LogInit()

	// init rpc
	FatalCheck(global.InitRpc(*rpc))
	monitor.C_Web_Script_Url = global.Config.RPC.ScriptUrl
	slog.Info("init rpc successed", "rpc", global.Config.RPC)

	// api prefix
	if global.Config.RPC.Chain != "" {
		global.Config.APIPrefix = fmt.Sprintf("/%s/api", global.Config.RPC.Chain)
	} else {
		global.Config.APIPrefix = "/api"
	}

	// db
	dbUser, dbPass := "", ""
	if *debug || *testnet {
		dbUser = "xxxxxx"
		dbPass = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	} else {
		dbUser = "xxxxxx"
		dbPass = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	}
	InitDB(*dbUrl, dbUser, dbPass, *dbName)
}

func InitDB(url, user, password, name string) {
	global.Config.DBUrl = url
	global.Config.DBName = name
	global.Config.DBUser = user
	global.Config.NodePath = *node
	global.Config.DBPassword = password
	FatalCheck(store.DBInit())
}

func InitLogic(ctx context.Context) {
	// websocket
	wss.Init(logic.StartKLinePush)

	// init cache
	FatalCheck(logic.InitCacheList(*hotList))
}

func InitCache() {
	// filesys
	FatalCheck(logic.InitFileSys())

	// init pairs
	FatalCheck(monitor.InitPairs())
}

func FatalCheck(params ...any) {
	if len(params) > 0 {
		if err, ok := params[len(params)-1].(error); ok && err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
