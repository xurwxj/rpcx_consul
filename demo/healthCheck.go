package demo

import (
	"fmt"
	"os"
	"time"

	"git.shining3d.com/cloud/util/types"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/xurwxj/rpcx_consul/registry"
)

var (
	registryServerHostPort string
	registryNameSpace      string
	registryAuth           string
	timeoutMs              uint64
	listenIntervalMs       uint64
	beatIntervalMs         int64
	serviceArr             []string
)
var s *server.Server

func initConfig() {
	registryAuth = "845e42bb-17cb-7a29-5ed8-11652dcb2a80"
	registryServerHostPort = "10.20.31.17:8500"
	registryNameSpace = "dev1"
}

func closeRegistry() {
	s.Close()
}
func InitRegistry() {
	initConfig()
	log.SetLogger(&types.RPCXServiceLogger{
		Logger: Log,
	})
	fmt.Println(time.Now().UTC().String(), "starting init registry....")
	s = server.NewServer()

	addConsulRegistryPlugins()

	if serviceArr == nil {
		serviceArr = []string{}
	}
	for _, sf := range serviceFuncs {
		serviceArr = append(serviceArr, sf.SFName)
		switch sf.SFType {
		case "func":
			s.RegisterFunction(sf.SFName, sf.SFCall, sf.SFMeta)
		case "class":
			s.RegisterName(sf.SFName, sf.SFCall, sf.SFMeta)
		}
	}
	addConsulCmuxPlugin()

	// go common.CallInitPermService()
	serviceAddr := "10.20.31.17:6974"
	fmt.Println(time.Now().UTC().String(), "starting registry server on ", serviceAddr)
	if err := s.Serve("tcp", serviceAddr); err != nil {
		fmt.Println("..", err, serviceAddr)
		os.Exit(1)
	}
}

// CloseRemoteRegistry close remote registry when closing
func CloseRemoteRegistry() {
	closeRegistry()

}

func addConsulRegistryPlugins() {
	fmt.Println("starting add consul registry plugin...")

	serviceAddr := "10.20.31.17:6974"

	if timeoutMs == 0 {
		timeoutMs = 3000
	}
	if listenIntervalMs == 0 {
		listenIntervalMs = 10000
	}
	if beatIntervalMs == 0 {
		beatIntervalMs = 5000
	}

	mode := "dev"

	r := &registry.ConsulRegisterPlugin{
		ServiceAddress:   "tcp@" + serviceAddr,
		ConsulServers:    registryServerHostPort,
		Token:            registryAuth,
		Datacenter:       registryNameSpace,
		Timeout:          timeoutMs / 1000,
		ListenIntervalMs: listenIntervalMs / 1000,
		BeatIntervalMs:   beatIntervalMs / 1000,
		Tags:             []string{mode},
		ENV:              mode,
		Log: &types.RPCXServiceLogger{
			Logger: Log,
		},
	}
	err := r.Start()
	if err != nil {
		fmt.Println("serviceAddr...", err, serviceAddr)
	}
	s.Plugins.Add(r)
}

func addConsulCmuxPlugin() {
	m := &registry.CmuxPlugin{
		Log: &types.RPCXServiceLogger{
			Logger: Log,
		},
		Services: &serviceArr,
	}
	m1 := &registry.CmuxPlugin1{
		CmuxPlugin: registry.CmuxPlugin{
			Log: &types.RPCXServiceLogger{
				Logger: Log,
			},
		},
	}
	s.Plugins.Add(m)
	s.Plugins.Add(m1)
}
