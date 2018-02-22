package main

import toml "github.com/pelletier/go-toml"

type Conf struct {
	Server ServerConf `toml:"server"`
	Auth   AuthConf   `toml:"auth"`
}

type ServerConf struct {
	Port            uint `toml:"port"`
	MaxHandlerCount uint `toml:"max_handler_count"`
}

type AuthConf struct {
	SecretKey string `toml:"secret_key"`
}

func NewConf(path string) (*Conf, error) {
	config, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Conf{}
	config.Unmarshal(conf)

	return conf, nil
}
