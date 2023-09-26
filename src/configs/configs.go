package configs

import (
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var vp *viper.Viper

func init() {
	vp = newViper()
}

func newViper() *viper.Viper {
	vp := viper.New()

	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vp.SetDefault("listen_port", "9101")

	vp.SetDefault("server_address", "127.0.0.1:9101")

	return vp
}

type Config interface {
	ServerConfig | ClientConfig
}

type Address string

type ServerConfig struct {
	ListenPort int `mapstructure:"listen_port"`
}

type ClientConfig struct {
	ServerAddress Address `mapstructure:"server_address"`
}

func decodeHookFunc() mapstructure.DecodeHookFuncType {
	return func(v, t reflect.Type, data interface{}) (interface{}, error) {

		if t == reflect.TypeOf(Address("")) {
			host, _, err := net.SplitHostPort(data.(string))
			if err != nil {
				return nil, err
			}
			if host == "" {
				return nil, errors.New("failed to parse address - empty host")
			}
		}

		return data, nil
	}
}

func load(cfgFile string) error {
	if cfgFile != "" {
		vp.SetConfigFile(cfgFile)
	} else {
		vp.AddConfigPath(".")
		vp.SetConfigType("yaml")
		vp.SetConfigName("messenger")
	}

	if err := vp.ReadInConfig(); err != nil {
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

	if err := vp.Unmarshal(&conf, viper.DecodeHook(decodeHookFunc())); err != nil {
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
