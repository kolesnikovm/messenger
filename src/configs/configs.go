package configs

import (
	"context"
	"errors"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const watchInterval = 5 * time.Second

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

	vp.SetDefault("consul.url", "localhost:8500")

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
	NewURL          []string      `mapstructure:"new_url"`
	Resharding      bool          `mapstructure:"resharding"`
	MaxConns        int32         `mapstructure:"max_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_connection_lifetime"`

	Changed chan struct{}
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
	vp *viper.Viper

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

	err := vp.AddRemoteProvider("consul", vp.GetString("consul.url"), "MESSENGER/CONFIG")
	if err != nil {
		log.Error().Err(err).Msg("unable to add remote provider")

		return vp, nil
	}

	vp.SetConfigType("json")

	err = vp.ReadRemoteConfig()
	if err != nil {
		log.Error().Err(err).Msg("unable to read remote config")
	}

	return vp, nil
}

func newConf[V Config](cfgFile string, conf *V) (*V, *viper.Viper, error) {
	vp, err := load(cfgFile)
	if err != nil {
		return conf, nil, err
	}

	err = vp.Unmarshal(conf, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToSliceHookFunc(","),
		decodeHookFunc(),
	)))
	if err != nil {
		return conf, nil, err
	}

	return conf, vp, nil
}

func NewClientConfig(cfgFile string) (*ClientConfig, error) {
	var conf ClientConfig

	clientConfig, _, err := newConf(cfgFile, &conf)

	return clientConfig, err
}

func NewServerConfig(cfgFile string) (*ServerConfig, error) {
	var conf ServerConfig

	serverConfig, vp, err := newConf(cfgFile, &conf)

	serverConfig.vp = vp

	serverConfig.Postgres.Changed = make(chan struct{})

	return serverConfig, err
}

func (c *ServerConfig) Watch(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(watchInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := c.vp.WatchRemoteConfig()
				if err != nil {
					log.Error().Err(err).Msg("unable to read remote config")
					continue
				}

				newConf := *c
				// nil slices in case of shrinking
				newConf.Postgres.URL = nil
				newConf.Postgres.NewURL = nil

				if err := c.vp.Unmarshal(&newConf); err != nil {
					log.Error().Err(err).Msg("unable to unmarshall remote config")
					continue
				}

				if cmp.Diff(c.Postgres, newConf.Postgres, cmpopts.SortSlices(func(a, b string) bool { return a < b })) != "" {
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c)), unsafe.Pointer(&newConf))
					c.Postgres.Changed <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
