package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg *Config

type Config struct {
	Mcs struct {
		Port      int    `yaml:"port"`
		Serverjar string `yaml:"server_jar"`
	}
	Bot struct {
		McsAddr               string   `yaml:"mcs_addr"`
		WhitelistedChannelIDS []string `yaml:"whitelisted_channel_ids"`
	}
}

func Load() error {
	cfgFile, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	Cfg = &Config{}

	d := yaml.NewDecoder(cfgFile)
	return d.Decode(Cfg)
}
