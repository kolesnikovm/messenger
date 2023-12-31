package configs

import (
	"errors"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func newViper() *viper.Viper {
	vp := viper.New()

	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vp.SetDefault("listen_port", "9101")
	vp.SetDefault("metrics_address", ":9102")

	vp.SetDefault("kafka.broker_list", "localhost:9094")

	vp.SetDefault("archiver.batch_size", 10)
	vp.SetDefault("archiver.num_senders", 4)
	vp.SetDefault("archiver.flush_interval", 1*time.Second)

	vp.SetDefault("postgres.url", "postgres://postgres:postgres@localhost:5432/messenger1,postgres://postgres:postgres@localhost:5432/messenger2")
	vp.SetDefault("postgres.max_connections", 10)
	vp.SetDefault("postgres.max_connection_lifetime", 10*time.Minute)

	vp.SetDefault("server_address", "127.0.0.1:9101")

	return vp
}

type Config interface {
	ServerConfig | ClientConfig
}

type Address struct {
	Host string
	Port int
}

type Postgres struct {
	URL             []string      `mapstructure:"url"`
	MaxConns        int32         `mapstructure:"max_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_connection_lifetime"`
}

type Archiver struct {
	BatchSize     int           `mapstructure:"batch_size"`
	NumSenders    int           `mapstructure:"num_senders"`
	FlushInterval time.Duration `mapstructure:"flush_interval"`
}

type Kafka struct {
	BrokerList []string `mapstructure:"broker_list"`
}

type ServerConfig struct {
	ListenPort     int      `mapstructure:"listen_port"`
	MetricsAddress string   `mapstructure:"metrics_address"`
	Kafka          Kafka    `mapstructure:"kafka"`
	Archiver       Archiver `mapstructure:"archiver"`
	Postgres       Postgres `mapstructure:"postgres"`
}

type ClientConfig struct {
	ServerAddress Address `mapstructure:"server_address"`
}

func decodeHookFunc() mapstructure.DecodeHookFuncType {
	return func(v, t reflect.Type, data interface{}) (interface{}, error) {

		if t == reflect.TypeOf(Address{}) {
			host, port, err := net.SplitHostPort(data.(string))
			if err != nil {
				return nil, err
			}
			if host == "" {
				return nil, errors.New("failed to parse address - empty host")
			}

			p, err := strconv.Atoi(port)
			if err != nil {
				return nil, err
			}

			return &Address{
				Host: host,
				Port: p,
			}, nil
		}

		return data, nil
	}
}

func load(cfgFile string) (*viper.Viper, error) {
	vp := newViper()

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
			log.Error().Err(err).Msg("")
		case *os.PathError:
			log.Error().Err(err).Msg("")
		default:
			return vp, err
		}
	}

	return vp, nil
}

func newConf[V Config](cfgFile string, conf V) (V, error) {
	vp, err := load(cfgFile)
	if err != nil {
		return conf, err
	}

	err = vp.Unmarshal(&conf, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToSliceHookFunc(","),
		decodeHookFunc(),
	)))
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func NewClientConfig(cfgFile string) (ClientConfig, error) {
	var conf ClientConfig
	return newConf(cfgFile, conf)
}

func NewServerConfig(cfgFile string) (ServerConfig, error) {
	var conf ServerConfig
	return newConf(cfgFile, conf)
}
