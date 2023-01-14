package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	Viper = viper.New()
	Env   = ""
)

const (
	Development = "dev"
	Production  = "prd"
	PreView     = "pre"
)

func init() {
	InitConf("")
}

func InitConf(confPath string) {
	Env = GetEnvironment()
	log.Println("start environment", Env)

	if confPath == "" {
		Viper.SetConfigName("app") //default app.toml
		confDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		confPath = filepath.Join(confDir, "../conf") + "/" + Env

		//Mac 本地goland编译
		if runtime.GOOS == "darwin" {
			_, fn, _, _ := runtime.Caller(0)
			confDir := filepath.Dir(fn)
			confPath = filepath.Join(confDir, "../../conf"+"/"+Env)
		}
		Viper.AddConfigPath(confPath)
	} else {
		Viper.SetConfigFile(confPath)
	}
	log.Println(confPath)

	if err := Viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal biz_error config file: %s", err))
	}

	log.Println(Env + " conf init success")
}

func GetEnvironment() string {
	path := "/home/sk/go-tx-chain/.deploy/service.su.txt"
	file, err := os.Open(path)
	if err != nil {
		return Development
	}

	data, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return Development
	}
	dataStr := string(data)

	switch {
	case strings.Contains(dataStr, "hnb-pre-v"):
		return PreView
	case strings.Contains(dataStr, "hnb-v"):
		return Production
	default:
		return Development
	}
}
