package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	GlobalAccounts = make(map[string]string)
	GlobalConfig   = new(Cfg)
)

type Cfg struct {
	*Server   `ini:"server"`
	*Ovpn     `ini:"ovpn"`
	*BaseAuth `ini:"baseauth"`
}

type Server struct {
	AppMode  string `ini:"app_mode"`
	HttpPort int    `ini:"http_port"`
	LogPath  string `ini:"log_path"`
}
type Ovpn struct {
	CertPath string `ini:"cert_path"`
}
type BaseAuth struct {
	Accounts []string `ini:"accounts"`
}

func parseAccounts(accounts []string, accountsMap map[string]string) (valid bool) {
	// ini []string split the account using ","
	for _, account := range accounts {
		// Split the username and password using "::"
		parts := strings.Split(account, "::")
		if len(parts) != 2 {
			log.Println("Invalid format for user value")
			return false
		}
		username := parts[0]
		password := parts[1]
		accountsMap[username] = password
	}
	return true
}

func Init() {
	configPath := flag.String("c", "openvpn-manager.ini", "config file path")
	flag.Parse()

	cfg, err := ini.Load(*configPath)
	if err != nil {
		log.Fatalf("Fail to load config file: %v", err)
	}

	GlobalConfig = &Cfg{
		Server: &Server{
			AppMode:  "production",
			HttpPort: 8080,
			LogPath:  "./",
		},
		Ovpn: &Ovpn{
			CertPath: "./ovpn",
		},
	}

	err = cfg.MapTo(GlobalConfig)
	if err != nil {
		fmt.Printf("ini map err: %v", err)
		os.Exit(1)
	}

	var valid bool
	valid = parseAccounts(GlobalConfig.Accounts, GlobalAccounts)
	if !valid {
		log.Printf("Invalid accounts config: %v, eg: user1::passwd1,user2::passwd2", GlobalConfig.Accounts)
		os.Exit(1)
	}

	// log.Printf("GlobalConfig: %s; GlobalAccounts: %s", GlobalConfig, GlobalAccounts)
}
