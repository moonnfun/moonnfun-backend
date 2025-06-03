package global

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	c_tgsession_timeout  = 86400
	c_websession_timeout = 3600
)

type gconf struct {
	// Server
	Log      string
	Port     int
	Mode     string
	WebPath  string
	HostURL  string
	NodePath string

	Testnet           bool
	WebSessionTimeout int

	// DB
	DBUrl      string
	DBName     string
	DBUser     string
	DBPassword string
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

	Config.WebSessionTimeout = c_websession_timeout
	cbuf, _ := json.MarshalIndent(Config, "", "	")

	fmt.Println(string(cbuf))
	return nil
}
