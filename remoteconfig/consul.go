package remoteconfig

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/smallnest/rpcx/log"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/rpcx_consul/client"
	"github.com/xurwxj/viper"
)

func initConsulConfig() {
	remoteF := base.CheckFileExistBackInfo("nacos/remote.json", true)
	cacheExist := false
	if remoteF != nil {
		remoteByte, err := os.ReadFile("nacos/remote.json")
		if err != nil {
			log.Errorf("initConsulConfig:ReadFile err %v", err)
		} else {
			err = viper.MergeConfig(bytes.NewReader(remoteByte))
			if err != nil {
				log.Errorf("initConsulConfig:MergeConfig err %v", err)
			} else {
				cacheExist = true
				go initConfig()
			}
		}
	}
	if !cacheExist {
		log.Infof("starting init consul config....")
		initConfig()
	}

}
func initConfig() {
	client, err := client.ConsulClient(configServerHostPort, configDataGroup, configAuth)
	if err != nil {
		log.Errorf("initConsulConfig configServerHostPort %v err %v", configServerHostPort, err)
		time.Sleep(10 * time.Second)
		initConsulConfig()
		return
	}
	// fmt.Println("client: ", client)
	c := client.KV()
	kvPath := fmt.Sprintf("%s/%s", configNameSpace, configDataID)
	kv, _, err := c.Get(kvPath, nil)
	if err != nil {
		log.Errorf("initConsulConfig kvPath =%v configServerHostPort=%v  err %v", kvPath, configServerHostPort, err)
		time.Sleep(10 * time.Second)
		initConsulConfig()
		return
	}
	// fmt.Println("kv: ", kv)
	if kv == nil {
		log.Errorf("initConsulConfig notFound  kvPath=%v configServerHostPort=%v ", kvPath, configServerHostPort)
		time.Sleep(10 * time.Second)
		initConsulConfig()
		return
	}
	MergeConfig(string(kv.Value))
}
