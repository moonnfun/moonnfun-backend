package global

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	C_Mode_Sync   = "sync"
	C_Mode_Server = "server"

	C_websession_timeout = 3 * 24 * 3600
)

type gconf struct {
	Debug bool

	Log       string
	Port      int
	Mode      string
	WebPath   string
	HostURL   string
	NodePath  string
	APIPrefix string

	RefreshSortListTimer   int
	RefreshTokenPriceTimer int

	Cmd        string
	Sync       bool
	Testnet    bool
	RPC        *RpcInfo
	LaunchMC   float64
	StartPrice float64

	ContractAddress   string
	WebSessionTimeout int

	// DB
	DBUrl      string
	DBName     string
	DBUser     string
	DBPassword string
}

type RpcInfo struct {
	ID            string
	Chain         string
	Http          string
	ScriptUrl     string
	Explorer      string
	Websocket     string
	WatchByEx     bool
	WatchByScan   bool
	WathchByLog   bool
	ValidChains   string
	WaitBlockTime int
	GetLogByFetch bool
	ScanPerTime   int
	FilterLogsMax int
	ChainAddrs    map[string]string

	ContractAddress    string
	WrapperAddress     string
	SwapRouterAddress  string
	TokenManageAddress string

	UserContractList string
}

var Config gconf

func Init(confPath string) error {
	if _, err := os.Stat(confPath); err != nil {
		return err
	}

	_, err := toml.DecodeFile(confPath, &Config)
	if err != nil {
		return err
	}
	cbuf, _ := json.MarshalIndent(Config, "", "	")

	fmt.Println(string(cbuf))
	return nil
}

func InitRpc(confPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &Config.RPC); err != nil {
		return err
	}
	if Config.RPC.WaitBlockTime == 0 {
		Config.RPC.WaitBlockTime = 100
	}
	Config.ContractAddress = Config.RPC.ContractAddress
	return nil
}
