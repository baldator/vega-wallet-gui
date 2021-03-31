package main

import (
	"fmt"
	"io/ioutil"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

const (
	configFile = "wallet-service-config.toml"
)

type Config struct {
	Host        string `json:"Host"`
	Level       string `json:"Level"`
	Port        int    `json:"Port"`
	RsaKey      string `json:"RsaKey"`
	TokenExpiry string `json:"TokenExpiry"`
	Console     struct {
		LocalPort int    `json:"LocalPort"`
		URL       string `json:"URL"`
	} `json:"Console"`
	Nodes struct {
		Hosts   []string `json:"Hosts"`
		Retries int      `json:"Retries"`
	} `json:"Nodes"`
}

type getConfigResponse struct {
	Path       string `json:"path"`
	PathExists bool   `json: "pathexist"`
	Config     Config `json: "config"`
}

func main() {
	var configResp getConfigResponse
	rootPath := fsutil.DefaultVegaDir()

	if ok, err := fsutil.PathExists(rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return
		}
		// create the folder
		if err := fsutil.EnsureDir(rootPath); err != nil {
			return
		}
	}

	if err := wallet.EnsureBaseFolder(rootPath); err != nil {
		return
	}
	configPath := rootPath + "/" + configFile
	configResp.Path = configPath
	configResp.PathExists = true

	if ok, err := fsutil.PathExists(rootPath); ok {
		fmt.Println("1")
		if _, ok := err.(*fsutil.PathNotFound); ok {
			fmt.Println("2")
			configResp.PathExists = false
		} else {
			fmt.Println("3")
			confFile, err := ioutil.ReadFile(configPath) // just pass the file name
			fmt.Println("4")
			if err != nil {
				fmt.Print(err)
			}
			config := Config{}
			fmt.Println("5")
			fmt.Printf("aaa: %v\n", confFile)
			toml.Unmarshal(confFile, &config)
			configResp.Config = config
		}
	} else {

		fmt.Println("nok")
	}

	fmt.Printf("%+v", configResp)
}

func initConfig(force bool, genRsaKey bool) error {
	rootPath := fsutil.DefaultVegaDir()
	if ok, err := fsutil.PathExists(rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(rootPath); err != nil {
			return fmt.Errorf("error creating root directory: %v", err)
		}
	}

	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	return wallet.GenConfig(log, rootPath, force, !genRsaKey)
}
