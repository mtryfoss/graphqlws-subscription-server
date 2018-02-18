package graphqlws_subscription_server

import toml "github.com/pelletier/go-toml"

type Conf struct {
	Port      uint   `toml:port`
	SecretKey string `toml:secret_key`
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
