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
	Grpc    grpcConf    `toml:"grpc"`
	Log     logConf     `toml:"log"`
	Routine routineConf `toml:"routine"`
	Store   storeConf   `toml:"store"`
}

type grpcConf struct {
	Addr string `toml:"addr"`
}

type logConf struct {
	FilePath string `toml:"file_path"`
}

type routineConf struct {
	MaxGoRoutineCount int `toml:"max_go_routine_count"`
}

type storeConf struct {
	Host []string `toml:"host"`
}
