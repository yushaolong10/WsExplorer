package main

import (
	"config"
	"flag"
	"fmt"
	"lib/logger"
	"server"
	"server/connection"
	"server/routine"
)

var (
	configFile = flag.String("c", "./conf/ws.toml.dev", "config file path")
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

	if err := connection.Init(config.Global.Conn.GroupCount, config.Global.Conn.MaxConnUserCount); err != nil {
		fmt.Printf("connection manager init err:%s", err.Error())
		return
	}
	fmt.Println("connection manager init success.")
	if err := routine.Init(config.Global.Routine.MaxGoRoutineCount); err != nil {
		fmt.Printf("routine init err:%s", err.Error())
		return
	}
	fmt.Println("routine init success.")
	//run
	fmt.Println("websocket server start...")
	logger.Info("websocket server start to work")
	server.Run()
}
