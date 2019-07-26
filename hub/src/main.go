package main

import (
	"config"
	"flag"
	"fmt"
	"lib/logger"
	"server"
)

var (
	configFile = flag.String("c", "./conf/hub.toml.dev", "config file path")
)

func main() {
	flag.Parse()
	//init
	if err := config.Init(*configFile); err != nil {
		fmt.Printf("config load err:%s", err.Error())
		return
	}
	fmt.Printf("config load success.\n%s\n", config.String())
	if err := logger.Init(config.Global.Log.FilePath, config.Global.Debug); err != nil {
		fmt.Printf("logger load err:%s", err.Error())
		return
	}
	fmt.Println("logger load success.")
	if err := server.Listen(config.Global.Grpc.Addr); err != nil {
		logger.Error("hub server start error. err:%s", err.Error())
	}
}
