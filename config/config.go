package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg *Config

type Config struct {
	Mcs struct {
		Port          int    `yaml:"port"`
		Serverjar     string `yaml:"server_jar"`
		EC2InstanceID string `yaml:"ec2_instance_id"`
	}
	Bot struct {
		McsAddr               string   `yaml:"mcs_addr"`
		McsPort               int      `yaml:"mcs_port"`
		WhitelistedChannelIDS []string `yaml:"whitelisted_channel_ids"`
	}
	Aws struct {
		Region       string `yaml:"region"`
		S3BucketName string `yaml:"s3_bucket_name"`
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
