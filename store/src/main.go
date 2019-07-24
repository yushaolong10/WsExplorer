package main

import (
	"config"
	"flag"
	"fmt"
	"logger"
	"server"
	"server/connection"
	"server/routine"
	"server/store"
)

var (
	configFile = flag.String("c", "./conf/store.toml.dev", "config file path")
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
	if err := routine.Init(config.Global.Routine.MaxGoRoutineCount); err != nil {
		fmt.Printf("routine init err:%s", err.Error())
		return
	}
	fmt.Println("routine init success.")
	if err := connection.Init(); err != nil {
		fmt.Printf("epoll init err:%s", err.Error())
		return
	}
	fmt.Println("epoll init success.")
	if err := store.Init(config.Global.Store.KeyMaxLen, config.Global.Store.ValMaxLen, config.Global.Store.SplitMaxDepth); err != nil {
		fmt.Printf("engine init err:%s", err.Error())
		return
	}
	fmt.Println("engine init success.")
	//run
	fmt.Println("store server start...")
	logger.Info("store server start to work")
	if err := server.Listen(config.Global.Tcp.Addr); err != nil {
		logger.Error("store server start error. err:%s", err.Error())
	}
}
