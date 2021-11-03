package remoteconfig

import (
	"bytes"
	"os"

	"github.com/smallnest/rpcx/log"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

var (
	configServerType        string
	configServerHostPort    string
	configServerContextPath string
	configNameSpace         string
	configDataID            string
	configDataGroup         string
	configAuth              string
	configCacheDir          string
	notLoadCacheAtStart     bool
	timeoutMs               uint64
	listenIntervalMs        uint64
)

// InitRemoteConfig init remote config and watch update, will override config value in config file
func InitRemoteConfig() {
	configServerType = viper.GetString("server.config.type")
	configServerHostPort = viper.GetString("server.config.hostPort")
	configServerContextPath = viper.GetString("server.config.contextPath")
	configNameSpace = viper.GetString("server.config.nameSpace")
	configDataID = viper.GetString("server.config.dataID")
	configDataGroup = viper.GetString("server.config.dataGroup")
	if configServerType == "" || configServerHostPort == "" || configNameSpace == "" || configDataID == "" || configDataGroup == "" {
		log.Errorf("please check config or cmd params configLost  func initRemoteConfig")
		os.Exit(1)
	}
	configAuth = viper.GetString("server.config.auth")
	configCacheDir = viper.GetString("server.config.cacheDir")
	notLoadCacheAtStart = viper.GetBool("server.config.notLoadCacheAtStart")
	timeoutMs = viper.GetUint64("server.config.timeoutMs")
	listenIntervalMs = viper.GetUint64("server.config.listenIntervalMs")
	switch configServerType {
	case "consul":
		initConsulConfig()
	}
}

func MergeConfig(data string) {
	if data == "" {
		return
	}
	base.CheckPathExistOrCreate("nacos")
	err := os.WriteFile("nacos/remote.json", []byte(data), os.ModePerm)
	if err != nil {
		log.Errorf("mergeConfig:WriteFile config=%v configServer=%v err=%v ", data, configServerType, err)
	}
	err = viper.MergeConfig(bytes.NewReader([]byte(data)))
	if err != nil {
		log.Errorf("mergeConfig config=%v configServer=%v err=%v ", data, configServerType, err)
	} else {
		log.Infof("mergeConfig done success configServer=%v", configServerType)
	}
}
