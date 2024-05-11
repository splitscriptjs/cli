package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Typescript bool     `toml:"typescript"`
	Main       string   `toml:"main"`
	Out        string   `toml:"out"`
	Dev        string   `toml:"dev"`
	Ignore     []string `toml:"ignore"`
}

func readConfig() (Config, error) {
	var conf Config
	_, err := toml.DecodeFile("splitscript.toml", &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
