package mcs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mcs struct {
		Port          int    `yaml:"port"`
		Serverjar     string `yaml:"server_jar"`
		EC2InstanceID string `yaml:"ec2_instance_id"`
	}
	Bot struct {
		WhitelistedChannelIDS []string `yaml:"whitelisted_channel_ids"`
	}
	Aws struct {
		Region       string `yaml:"region"`
		S3BucketName string `yaml:"s3_bucket_name"`
	}
}

func LoadConfig() (*Config, error) {
	cfgFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	cfg := &Config{}

	d := yaml.NewDecoder(cfgFile)
	err = d.Decode(cfg)
	return cfg, err
}
