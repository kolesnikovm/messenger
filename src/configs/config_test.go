package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	goodConf = `listen_port: "9102"
server_address: "127.0.0.1:9102"`

	badConf = `badConf`
)

func createConfig(config string) error {
	return os.WriteFile("messenger.tmp.yaml", []byte(config), 0644)
}

func cleanupConfig() error {
	return os.Remove("messenger.tmp.yaml")
}

func TestNewClientConfig(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func()
		teardown  func()
		cfgFile   string
		exrResult Address
	}{
		{
			name:      "default_config",
			setup:     func() {},
			teardown:  func() {},
			cfgFile:   "",
			exrResult: Address{Host: "127.0.0.1", Port: 9101},
		},
		{
			name: "file_config",
			setup: func() {
				if err := createConfig(goodConf); err != nil {
					t.Fatal(err)
				}
			},
			teardown: func() {
				if err := cleanupConfig(); err != nil {
					t.Fatal(err)
				}
			},
			cfgFile:   "./messenger.tmp.yaml",
			exrResult: Address{Host: "127.0.0.1", Port: 9102},
		},
		{
			name:      "env_config",
			setup:     func() { os.Setenv("SERVER_ADDRESS", "127.0.0.1:9103") },
			teardown:  func() { os.Unsetenv("SERVER_ADDRESS") },
			exrResult: Address{Host: "127.0.0.1", Port: 9103},
		},
		{
			name:      "not_found",
			setup:     func() { os.Setenv("SERVER_ADDRESS", "127.0.0.1:9103") },
			teardown:  func() { os.Unsetenv("SERVER_ADDRESS") },
			cfgFile:   "not/exists.yaml",
			exrResult: Address{Host: "127.0.0.1", Port: 9103},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()
			defer testCase.teardown()

			config, err := NewClientConfig(testCase.cfgFile)

			require.NoError(t, err)
			require.Equal(t, testCase.exrResult, config.ServerAddress)
		})
	}
}

func TestNewClientConfigError(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func()
		teardown func()
		cfgFile  string
	}{
		{
			name:     "wrong_extention",
			setup:    func() {},
			teardown: func() {},
			cfgFile:  "/not/exists",
		},
		{
			name: "wrong_format",
			setup: func() {
				if err := createConfig(badConf); err != nil {
					t.Fatal(err)
				}
			},
			teardown: func() {
				if err := cleanupConfig(); err != nil {
					t.Fatal(err)
				}
			},
			cfgFile: "./messenger.tmp.yaml",
		},
		{
			name:     "wrong_value",
			setup:    func() { os.Setenv("SERVER_ADDRESS", "text") },
			teardown: func() { os.Unsetenv("SERVER_ADDRESS") },
		},
		{
			name:     "wrong_host",
			setup:    func() { os.Setenv("SERVER_ADDRESS", ":9101") },
			teardown: func() { os.Unsetenv("SERVER_ADDRESS") },
		},
		{
			name:     "wrong_port",
			setup:    func() { os.Setenv("SERVER_ADDRESS", "host:port") },
			teardown: func() { os.Unsetenv("SERVER_ADDRESS") },
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()
			defer testCase.teardown()

			_, err := NewClientConfig(testCase.cfgFile)

			require.Error(t, err)
		})
	}
}

func TestNewServerConfig(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func()
		teardown  func()
		cfgFile   string
		exrResult int
	}{
		{
			name:      "default_config",
			setup:     func() {},
			teardown:  func() {},
			cfgFile:   "",
			exrResult: 9101,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()
			defer testCase.teardown()

			config, err := NewServerConfig(testCase.cfgFile)

			require.NoError(t, err)
			require.Equal(t, testCase.exrResult, config.ListenPort)
		})
	}
}

func TestNewServerConfigError(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func()
		teardown func()
		cfgFile  string
	}{
		{
			name:     "wrong_value",
			setup:    func() { os.Setenv("LISTEN_PORT", "text") },
			teardown: func() { os.Unsetenv("LISTEN_PORT") },
			cfgFile:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()
			defer testCase.teardown()

			_, err := NewServerConfig(testCase.cfgFile)

			require.Error(t, err)
		})
	}
}
