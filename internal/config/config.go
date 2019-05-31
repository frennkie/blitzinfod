package config

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	defaultBlitzdDirName    = "blitzd"
	defaultBlitzdConfigName = "config.toml"

	defaultCfgFile      = "/etc/blitzd.toml"
	defaultCfgFileWin32 = "C:\\blitzd.toml" // Win32 mostly used for development

	defaultBlitzdAlias = "MyBlitz42"

	defaultEnvPrefix = "BLITZD"

	defaultTLSServerCaCertFilename = "blitzd_ca.crt"
	defaultTLSServerCertFilename   = "blitzd_server.crt"
	defaultTLSServerKeyFilename    = "blitzd_server.key"

	defaultTLSClientCaCertFilename = "blitzd_ca.crt"
	defaultTLSClientCertFilename   = "blitzd_client.crt"
	defaultTLSClientKeyFilename    = "blitzd_client.key"

	defaultHttpPort     = "38080"
	defaultHttpHostPort = "localhost:" + defaultHttpPort

	defaultHttpsPort     = "38443"
	defaultHttpsHostPort = "localhost:" + defaultHttpsPort

	defaultRPCPort     = "39735"
	defaultRPCHostPort = "localhost:" + defaultRPCPort
)

var (
	DefaultBlitzdDir = btcutil.AppDataDir(defaultBlitzdDirName, false)
	BlitzdDir        string

	// ToDo check
	// maxMsgRecvSize is the largest message our client will receive. We
	// set this to ~50Mb atm.
	//maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 50)

)

func setDefaults() {
	log.Printf("Settings Defaults...")

	viper.SetDefault("blitzdDir", BlitzdDir)

	viper.SetDefault("customCfgPath", "")
	viper.SetDefault("defaultCfgPath", "")

	viper.SetDefault("alias", defaultBlitzdAlias)

	viper.SetDefault("server.cacert", filepath.Join(BlitzdDir, defaultTLSServerCaCertFilename))
	viper.SetDefault("server.tlscert", filepath.Join(BlitzdDir, defaultTLSServerCertFilename))
	viper.SetDefault("server.tlskey", filepath.Join(BlitzdDir, defaultTLSServerKeyFilename))
	viper.SetDefault("client.cacert", filepath.Join(BlitzdDir, defaultTLSClientCaCertFilename))
	viper.SetDefault("client.tlscert", filepath.Join(BlitzdDir, defaultTLSClientCertFilename))
	viper.SetDefault("client.tlskey", filepath.Join(BlitzdDir, defaultTLSClientKeyFilename))

	viper.SetDefault("restHostPort", defaultHttpHostPort)
	viper.SetDefault("rpcHostPort", defaultRPCHostPort)

	viper.SetDefault("server.http.enabled", true)
	viper.SetDefault("server.http.port", 30080)

	viper.SetDefault("server.https.enabled", true)
	viper.SetDefault("server.https.port", 30080)

	viper.SetDefault("server.rpc.enabled", true)
	viper.SetDefault("server.rpc.port", defaultRPCPort)

}

func InitConfig() {
	// First set default values.
	// Then read default (/etc/blitzd.toml|C:\blitzd.toml)
	// Then - if it exists - merge any settings from file "blitzd.toml" in home directory
	// Then load env
	// Finally activate the Config Watcher
	setDefaults()

	// read default
	if runtime.GOOS == "windows" {
		viper.Set("defaultCfgPath", defaultCfgFileWin32)
		viper.SetConfigFile(filepath.FromSlash(defaultCfgFileWin32))
	} else {
		viper.Set("defaultCfgPath", defaultCfgFile)
		viper.SetConfigFile(filepath.FromSlash(defaultCfgFile))
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// check custom config
	customCfgPath := filepath.FromSlash(filepath.Join(BlitzdDir, defaultBlitzdConfigName))
	if _, err := os.Stat(customCfgPath); os.IsNotExist(err) {
		log.Printf("custom config file does not exist - skipping: %s", customCfgPath)
	} else {
		viper.Set("customCfgPath", customCfgPath)
		viper.SetConfigFile(customCfgPath)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Println("Can't read config for merge:", err)
			os.Exit(1)
		}

		log.Printf("Merged config file: %s", customCfgPath)
	}

	// load env
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// config watcher
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
	})

	// store copy of parsed/merged config
	_ = viper.WriteConfigAs(filepath.Join(BlitzdDir, "saved.toml"))

}