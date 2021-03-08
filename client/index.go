package client

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
)

// ConsulClient get consul client by config params
func ConsulClient(serverStrs, datacenter, token string) (client *api.Client, err error) {
	scStrs := strings.Split(strings.TrimSpace(serverStrs), ",")
	if len(scStrs) < 1 {
		err = fmt.Errorf("noServerConfig")
		return
	}
	conf := api.DefaultConfig()
	conf.Datacenter = datacenter
	conf.Token = token

	for _, sc := range scStrs {
		conf.Address = strings.TrimSpace(sc)
		client, err = api.NewClient(conf)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("noAvaiableServer")
	return
}
