package configs

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("server.listenAddress", "0.0.0.0:9101")

	viper.SetDefault("client.serverAddress", "127.0.0.1:9101")
}

type Server struct {
	ListenAddress string `mapstructure:"listenAddress"`
}

type Client struct {
	ServerAddress string `mapstructure:"serverAddress"`
}

type Config struct {
	Server Server `mapstructure:"server"`
	Client Client `mapstructure:"client"`
}

func Load(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("messenger")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func New() (Config, error) {
	var conf Config

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}
