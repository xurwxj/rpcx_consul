package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/util"
)

// ConsulRegisterPlugin implements consul registry.
type ConsulRegisterPlugin struct {
	// service address, for example, tcp@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string
	// Consul client config
	ClientConfig *api.Config

	Timeout          uint64
	ListenIntervalMs uint64
	BeatIntervalMs   int64

	// Registered services
	Services []string
	Tags     []string
	ENV      string

	namingClient *api.Agent
	Client       *api.Client

	dying chan struct{}
	done  chan struct{}
}

// Start starts to connect consul cluster
func (p *ConsulRegisterPlugin) Start() error {
	if p.done == nil {
		p.done = make(chan struct{})
	}
	if p.dying == nil {
		p.dying = make(chan struct{})
	}

	client, err := api.NewClient(p.ClientConfig)
	if err != nil {
		return err
	}
	p.Client = client
	p.namingClient = client.Agent()

	return nil
}

// Stop unregister all services.
func (p *ConsulRegisterPlugin) Stop() error {

	for _, name := range p.Services {
		err := p.namingClient.ServiceDeregister(name)
		if err != nil {
			log.Errorf("Stop:ServiceDeregister:%v, with name: %s", err, name)
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

	inst := &api.AgentServiceRegistration{
		Name:    name,
		ID:      fmt.Sprintf("%s_%s_%v", name, ip, port),
		Port:    port,
		Tags:    p.Tags,
		Address: ip,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%v", ip, port),
			Timeout:                        fmt.Sprintf("%vs", p.Timeout),
			Interval:                       fmt.Sprintf("%vs", p.BeatIntervalMs),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%vs", p.ListenIntervalMs),
		},
	}

	err = p.namingClient.ServiceRegister(inst)
	if err != nil {
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
	err = p.namingClient.ServiceDeregister(name)
	if err != nil {
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
