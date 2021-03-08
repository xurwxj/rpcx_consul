package discovery

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/util"
)

// ConsulDiscovery is a consul service discovery.
// It always returns the registered servers in consul.
type ConsulDiscovery struct {
	servicePath string
	env         string
	tag         string
	// consul client config
	clientConfig *api.Config

	Client *api.Client

	pairsMu sync.RWMutex
	pairs   []*client.KVPair
	chans   []chan []*client.KVPair
	mu      sync.Mutex

	filter                  client.ServiceDiscoveryFilter
	RetriesAfterWatchFailed int

	stopCh chan struct{}
}

// NewConsulDiscovery returns a new ConsulDiscovery.
func NewConsulDiscovery(servicePath, env, tag string, clientConfig *api.Config) (client.ServiceDiscovery, error) {
	d := &ConsulDiscovery{
		servicePath:  servicePath,
		env:          env,
		tag:          tag,
		clientConfig: clientConfig,
	}

	client, err := api.NewClient(d.clientConfig)
	if err != nil {
		return nil, err
	}
	d.Client = client

	d.fetch()
	go d.watch()
	return d, nil
}

// NewConsulDiscoveryWithClient returns a new ConsulDiscovery with *api.Config
func NewConsulDiscoveryWithClient(servicePath, env, tag string, client *api.Client) client.ServiceDiscovery {
	d := &ConsulDiscovery{
		servicePath: servicePath,
		env:         env,
		tag:         tag,
	}

	d.Client = client

	d.fetch()
	go d.watch()
	return d
}

func (d *ConsulDiscovery) fetch() {
	// var lastIndex uint64
	// services, metainfo, err := d.Client.Health().Service(d.servicePath, d.tag, true, &api.QueryOptions{
	// 	WaitIndex: lastIndex,
	// })
	_, services, err := d.Client.Agent().AgentHealthServiceByName(d.servicePath)
	if err != nil {
		log.Errorf("failed to get service %s: %v", d.servicePath, err)
		return
	}
	// lastIndex = metainfo.LastIndex
	pairs := make([]*client.KVPair, 0, len(services))
	for _, inst := range services {
		if FindInStringSlice(inst.Service.Tags, d.env) && inst.AggregatedStatus == "passing" {
			network := inst.Service.Meta["network"]
			ip := inst.Service.Address
			port := inst.Service.Port
			key := fmt.Sprintf("%s@%s:%d", network, ip, port)
			pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Service.Meta)}
			if d.filter != nil && !d.filter(pair) {
				continue
			}
			pairs = append(pairs, pair)
		}
	}

	d.pairsMu.Lock()
	d.pairs = pairs
	d.pairsMu.Unlock()
}

// FindInStringSlice find string in slice
func FindInStringSlice(s []string, t string) bool {
	for _, v := range s {
		if strings.TrimSpace(v) == strings.TrimSpace(t) {
			return true
		}
	}
	return false
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d *ConsulDiscovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	return NewConsulDiscovery(servicePath, d.env, d.tag, d.clientConfig)
}

// SetFilter sets the filer.
func (d *ConsulDiscovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *ConsulDiscovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

// WatchService returns a nil chan.
func (d *ConsulDiscovery) WatchService() chan []*client.KVPair {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*client.KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

// RemoveWatcher remove wather
func (d *ConsulDiscovery) RemoveWatcher(ch chan []*client.KVPair) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*client.KVPair
	for _, c := range d.chans {
		if c == ch {
			continue
		}

		chans = append(chans, c)
	}

	d.chans = chans
}

func (d *ConsulDiscovery) watch() {
	// param := &vo.SubscribeParam{
	// 	ServiceName: d.servicePath,
	// 	Clusters:    []string{d.Cluster},
	// 	SubscribeCallback: func(services []model.SubscribeService, err error) {
	// 		pairs := make([]*client.KVPair, 0, len(services))
	// 		for _, inst := range services {
	// 			network := inst.Metadata["network"]
	// 			ip := inst.Ip
	// 			port := inst.Port
	// 			key := fmt.Sprintf("%s@%s:%d", network, ip, port)
	// 			pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Metadata)}
	// 			if d.filter != nil && !d.filter(pair) {
	// 				continue
	// 			}
	// 			pairs = append(pairs, pair)
	// 		}
	// 		d.pairsMu.Lock()
	// 		d.pairs = pairs
	// 		d.pairsMu.Unlock()

	// 		d.mu.Lock()
	// 		for _, ch := range d.chans {
	// 			ch := ch
	// 			go func() {
	// 				defer func() {
	// 					recover()
	// 				}()
	// 				select {
	// 				case ch <- d.pairs:
	// 				case <-time.After(time.Minute):
	// 					log.Warn("chan is full and new change has been dropped")
	// 				}
	// 			}()
	// 		}
	// 		d.mu.Unlock()
	// 	},
	// }

	// err := d.namingClient.Subscribe(param)
	// // if failed to Subscribe, retry
	// if err != nil {
	// 	var tempDelay time.Duration
	// 	retry := d.RetriesAfterWatchFailed
	// 	for d.RetriesAfterWatchFailed < 0 || retry >= 0 {
	// 		err := d.namingClient.Subscribe(param)
	// 		if err != nil {
	// 			if d.RetriesAfterWatchFailed > 0 {
	// 				retry--
	// 			}
	// 			if tempDelay == 0 {
	// 				tempDelay = 1 * time.Second
	// 			} else {
	// 				tempDelay *= 2
	// 			}
	// 			if max := 30 * time.Second; tempDelay > max {
	// 				tempDelay = max
	// 			}
	// 			log.Warnf("can not subscribe (with retry %d, sleep %v): %s: %v", retry, tempDelay, d.servicePath, err)
	// 			time.Sleep(tempDelay)
	// 			continue
	// 		}
	// 		break
	// 	}
	// }
}

// Close close service client
func (d *ConsulDiscovery) Close() {
	close(d.stopCh)
}
