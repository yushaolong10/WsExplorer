package main

import (
	"config"
	"flag"
	"fmt"
	"logger"
	"server"
	"server/connection"
	"server/routine"
)

var (
	configFile = flag.String("c", "./conf/websocket.toml.dev", "config file path")
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
	fmt.Printf("logger load success.")

	if err := connection.Init(config.Global.Conn.GroupCount, config.Global.Conn.MaxConnUserCount); err != nil {
		fmt.Printf("connection manager init err:%s", err.Error())
		return
	}
	fmt.Printf("connection manager init success.")
	if err := routine.Init(config.Global.Routine.MaxGoRoutineCount); err != nil {
		fmt.Printf("routine init err:%s", err.Error())
		return
	}
	fmt.Printf("routine init success.")
	//run
	server.Run()
}
