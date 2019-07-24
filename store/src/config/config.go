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

var Global storeConfig

type storeConfig struct {
	Cluster string      `toml:"cluster"`
	Debug   bool        `toml:"debug"`
	Env     string      `toml:"env"`
	Tcp     tcpConf     `toml:"tcp"`
	Log     logConf     `toml:"log"`
	Routine routineConf `toml:"routine"`
	Store   storeConf   `toml:"store"`
}

type tcpConf struct {
	Addr string `toml:"addr"`
}

type logConf struct {
	FilePath string `toml:"file_path"`
}

type routineConf struct {
	MaxGoRoutineCount int `toml:"max_go_routine_count"`
}

type storeConf struct {
	KeyMaxLen     int `toml:"key_max_len"`
	ValMaxLen     int `toml:"val_max_len"`
	SplitMaxDepth int `toml:"split_max_depth"`
}
