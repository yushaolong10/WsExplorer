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
	Grpc    grpcConf    `toml:"grpc"`
	Store   storeConf   `toml:"store"`
	Hub     hubConf     `toml:"hub"`
	Log     logConf     `toml:"log"`
	Conn    connConf    `toml:"connection"`
	Routine routineConf `toml:"routine"`
}

type httpConf struct {
	Addr string `toml:"addr"`
}

type grpcConf struct {
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

type storeConf struct {
	Host []string `toml:"host"`
	Pool poolConf `toml:"pool"`
}

type hubConf struct {
	Host []string `toml:"host"`
	Pool poolConf `toml:"pool"`
}

type poolConf struct {
	MinOpen     int `toml:"min_open"`
	MaxOpen     int `toml:"max_open"`
	MaxLifeTime int `toml:"max_life_time"`
	Timeout     int `toml:"timeout"`
}
