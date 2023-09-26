package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("listen_address", "0.0.0.0:9101")

	viper.SetDefault("server_address", "127.0.0.1:9101")
}

type Config interface {
	ServerConfig | ClientConfig
}

type ServerConfig struct {
	ListenAddress string `mapstructure:"listen_address"`
}

type ClientConfig struct {
	ServerAddress string `mapstructure:"server_address"`
}

func load(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("messenger")
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Printf("%s\n", err)
		case *os.PathError:
			fmt.Printf("%s\n", err)
		default:
			return err
		}
	}

	return nil
}

func new[V Config](cfgFile string, conf V) (V, error) {
	if err := load(cfgFile); err != nil {
		return conf, err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func NewClientConfig(cfgFile string) (ClientConfig, error) {
	var conf ClientConfig
	return new(cfgFile, conf)
}

func NewServerConfig(cfgFile string) (ServerConfig, error) {
	var conf ServerConfig
	return new(cfgFile, conf)
}
