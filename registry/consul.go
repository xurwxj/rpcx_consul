package registry

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/util"
	"github.com/soheilhy/cmux"
)

// ConsulRegisterPlugin implements consul registry.
type ConsulRegisterPlugin struct {
	// service address, for example, tcp@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string

	HealthType string

	HttpHealthPort string
	// consul client config Datacenter
	Datacenter string
	// consul server token
	Token string
	// consul server address, multi seperated by comma
	ConsulServers string
	// Consul client config
	ClientConfig *api.Config

	Timeout          uint64
	ListenIntervalMs uint64
	BeatIntervalMs   int64
	Log              log.Logger

	// Registered services
	Services []string
	Tags     []string
	ENV      string

	namingClients []*api.Agent
	Clients       []*api.Client

	dying chan struct{}
	done  chan struct{}
}

// Start starts to connect consul cluster
func (p *ConsulRegisterPlugin) Start(l net.Listener) error {
	log.SetLogger(p.Log)
	if p.done == nil {
		p.done = make(chan struct{})
	}
	if p.dying == nil {
		p.dying = make(chan struct{})
	}
	if p.HealthType == "" {
		p.HealthType = "tcp"
	}
	if p.HealthType == "http" && p.HttpHealthPort != "" {
		// _, ip, _, err := util.ParseRpcxAddress(p.ServiceAddress)
		// if err != nil {
		// 	return err
		// }
		// server := &http.Server{
		// 	Addr:         fmt.Sprintf("%s:%s", ip, p.HttpHealthPort),
		// 	ReadTimeout:  10 * time.Second,
		// 	WriteTimeout: 10 * time.Second,
		// }
		// mux := http.NewServeMux()
		m := cmux.New(l)
		fmt.Println("m ....", m)
		httpL := m.Match(cmux.HTTP1Fast())
		httpS := &http.Server{}
		go httpS.Serve(httpL)
		err := m.Serve()
		fmt.Println("m serve....", err)
		// mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// 	serviceName := r.URL.Query().Get("service")
		// 	if serviceName != "" && checkServiceExist(serviceName, p) {
		// 		resByte, err := base.GetByteArrayFromInterface(map[string]interface{}{
		// 			"status": "UP",
		// 			"application": map[string]string{
		// 				"status": "UP",
		// 			},
		// 		})
		// 		if err == nil {
		// 			_, err = w.Write(resByte)
		// 			if err == nil {
		// 				return
		// 			}
		// 		}
		// 	}
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("fail"))
		// })
		// server = mux
		// server.Handler = mux
		// go func() {
		// 	fmt.Println("http start...")
		// 	err := server.ListenAndServe()
		// 	if err != nil {
		// 		fmt.Println("http: ", err)
		// 	}
		// }()
	}
	scStrs := strings.Split(strings.TrimSpace(p.ConsulServers), ",")
	if len(scStrs) < 1 {
		return fmt.Errorf("noServerConfig")
	}
	conf := api.DefaultConfig()
	conf.Datacenter = p.Datacenter
	conf.Token = p.Token
	for _, cs := range scStrs {
		conf.Address = strings.TrimSpace(cs)
		client, err := api.NewClient(conf)
		if err != nil {
			log.Errorf("ConsulRegisterPlugin:NewClient:err: %v on server: %s", err, cs)
			continue
		}
		log.Debugf("ConsulRegisterPlugin:NewClient:server: %s", cs)
		p.namingClients = append(p.namingClients, client.Agent())
		p.Clients = append(p.Clients, client)
	}
	// client.Catalog().Register(&api.CatalogRegistration{})

	return nil
}

func checkServiceExist(serviceName string, p *ConsulRegisterPlugin) (rs bool) {
	// //test
	rs = true
	if serviceName == "com.shining3d.sm.log" {
		if rand.Intn(2) == 1 {
			rs = false
		}
	}
	return
	// for _, s := range p.Services {
	// 	if s == serviceName {
	// 		rs = true
	// 	}
	// }
	// return
}

// Stop unregister all services.
func (p *ConsulRegisterPlugin) Stop() error {

	for _, name := range p.Services {
		for _, nc := range p.namingClients {
			err := nc.ServiceDeregister(name)
			if err != nil {
				log.Errorf("Stop:ServiceDeregister:%v, with name: %s", err, name)
			}
		}
	}

	close(p.dying)
	<-p.done

	return nil
}

// Register handles registering event.
// this service is registered at BASE/serviceName/thisIpAddress node
func (p *ConsulRegisterPlugin) Register(name string, rcvr interface{}, metadata string) (err error) {
	if strings.TrimSpace(name) == "" {
		return errors.New("serviceEmpty")
	}

	network, ip, port, err := util.ParseRpcxAddress(p.ServiceAddress)
	if err != nil {
		return err
	}

	meta := util.ConvertMeta2Map(metadata)
	meta["network"] = network
	meta["env"] = p.ENV

	healthCheck := api.AgentServiceCheck{
		Timeout:                        fmt.Sprintf("%vs", p.Timeout),
		Interval:                       fmt.Sprintf("%vs", p.BeatIntervalMs),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%vs", p.ListenIntervalMs),
	}
	if p.HealthType == "tcp" {
		healthCheck.TCP = fmt.Sprintf("%s:%v", ip, port)
	}
	if p.HealthType == "http" && p.HttpHealthPort != "" {
		healthCheck.HTTP = fmt.Sprintf("http://%s:%s/health?service=%s", ip, p.HttpHealthPort, name)
	}

	inst := &api.AgentServiceRegistration{
		Name:    name,
		ID:      fmt.Sprintf("%s_%s_%v", name, ip, port),
		Port:    port,
		Tags:    p.Tags,
		Address: ip,
		Meta:    meta,
		Check:   &healthCheck,
	}
	done := false
	for _, nc := range p.namingClients {
		err = nc.ServiceRegister(inst)
		ncn, _ := nc.NodeName()
		if err == nil {
			log.Debugf("ConsulRegisterPlugin:NewClient:server: %s:service:%s", ncn, inst.Name)
			done = true
		} else {
			log.Errorf("ConsulRegisterPlugin:NewClient:server: %s:service:%s:err:%v", ncn, inst.Name, err)
		}
	}
	if err != nil && !done {
		return err
	}

	p.Services = append(p.Services, name)

	return
}

// RegisterFunction register function
func (p *ConsulRegisterPlugin) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return p.Register(serviceName, fn, metadata)
}

// Unregister Unregister service
func (p *ConsulRegisterPlugin) Unregister(name string) (err error) {
	if len(p.Services) == 0 {
		return nil
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("Unregister service `name` can't be empty")
	}

	// _, ip, port, err := util.ParseRpcxAddress(p.ServiceAddress)
	// if err != nil {
	// 	return err
	// }
	done := false
	for _, nc := range p.namingClients {
		err = nc.ServiceDeregister(name)
		if err == nil {
			done = true
		}
	}
	if err != nil && !done {
		return err
	}

	var services = make([]string, 0, len(p.Services)-1)
	for _, s := range p.Services {
		if s != name {
			services = append(services, s)
		}
	}
	p.Services = services

	return nil
}
