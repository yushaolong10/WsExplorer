package config

import (
	"github.com/BurntSushi/toml"
	"lib/json"
)

func Init(filePth string) error {
	_, err := toml.DecodeFile(filePth, &Global)
	return err
}

func String() string {
	b, _ := json.Marshal(Global)
	return string(b)
}

var Global wsConfig

type wsConfig struct {
	Cluster string      `toml:"cluster"`
	Debug   bool        `toml:"debug"`
	Env     string      `toml:"env"`
	Http    httpConf    `toml:"http"`
	Log     logConf     `toml:"log"`
	Conn    connConf    `toml:"connection"`
	Routine routineConf `toml:"routine"`
}

type httpConf struct {
	Addr string `toml:"addr"`
}

type logConf struct {
	FilePath string `toml:"file_path"`
}

type connConf struct {
	GroupCount       int `toml:"group_count"`
	MaxConnUserCount int `toml:"max_conn_user_count"`
}

type routineConf struct {
	MaxGoRoutineCount int `toml:"max_go_routine_count"`
}
