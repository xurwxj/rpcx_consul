package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/xurwxj/rpcx_consul/discovery"
)

func main() {
	conf := api.DefaultConfig()
	conf.Address = "10.10.1.58:8500"
	conf.Token = "1c67c216-a4ec-2235-da71-ad5d3f970280"
	client, err := api.NewClient(conf)
	if err != nil {
		fmt.Println("client err: ", err)
	}
	fmt.Println(time.Now())
	serv, err := client.Agent().ServicesWithFilter("dev in Tags")
	if err != nil {
		fmt.Println("client Service err: ", err)
	}
	fmt.Println("serv: ", serv)
	for _, m := range serv {
		if m.Service == "com.shining3d.app.auth.audit" {
			fmt.Println("m: ", m.Service)
			fmt.Println("m: ", m.Address)
			fmt.Println("m: ", m.Port)
		}
	}
	fmt.Println(time.Now())

	var lastIndex uint64
	ssss, smi, err := client.Catalog().Service("com.shining3d.app.auth.audit", "dev", &api.QueryOptions{
		WaitIndex: lastIndex,
	})
	if err != nil {
		fmt.Println("client Service err: ", err)
	}
	fmt.Println("ssss: ", ssss)
	fmt.Println("smi: ", smi)
	fmt.Println(time.Now())
	rs, siii, err := client.Agent().AgentHealthServiceByName("com.shining3d.app.auth.audit")
	if err != nil {
		fmt.Println("client Service err: ", err)
	}
	fmt.Println("rs: ", rs)
	fmt.Println("siii: ", siii)
	for _, k := range siii {
		if discovery.FindInStringSlice(k.Service.Tags, "dev") {
			fmt.Println("k", k.Service.Service)
			fmt.Println("k", k.Service.Address)
			fmt.Println("k", k.Service.Port)
		}
	}
	fmt.Println(time.Now())
	// services, metainfo, err := client.Health().Service("com.shining3d.app.auth.audit", "dev", true, &api.QueryOptions{
	// 	WaitIndex: lastIndex,
	// })
	// if err != nil {
	// 	fmt.Println("client Service err: ", err)
	// }
	// lastIndex = metainfo.LastIndex
	// for _, service := range services {
	// 	fmt.Println("address: ", service.Service.Address)
	// }
}
