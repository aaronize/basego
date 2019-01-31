package utils

import (
	"configcenter/src/framework/core/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Listen 		string
	Logger 		*Logger

	RDB 		*RDB
	Kafka 		*Kafka
	Redis 		*Redis

	Debug 		bool
}

type Redis struct {
	Cluster 	[]string
	PoolSize 	int
	Host 		string
	Password 	string
	DB 			int
}

type Kafka struct {
	Hosts 	[]string
	Topic  	[]string
}

type RDB struct {
	Database 	string
	Host 		string
	Port 		string
	User 		string
	Password  	string
}

type Logger struct {
	Path 		string
	MaxSize 	int
	MaxBackups 	int
	MaxAge 		int
}

// TODO 设置默认值
func ParseConfig(path string) (*Config, error) {
	var config *Config

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("load %s file error. ", path))
	}

	if len(content) == 0 {
		return nil, errors.New(fmt.Sprintf("file %s is empty", path))
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, errors.New(fmt.Sprintf("parse %s file error. ", path))
	}

	return config, nil
}