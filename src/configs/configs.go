package configs

import (
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.listen_address", "0.0.0.0:9101")

	viper.SetDefault("client.server_address", "127.0.0.1:9101")
}

type Server struct {
	ListenAddress string `mapstructure:"listen_address"`
}

type Client struct {
	ServerAddress string `mapstructure:"server_address"`
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
