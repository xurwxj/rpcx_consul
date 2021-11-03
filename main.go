package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xurwxj/rpcx_consul/demo"
	"github.com/xurwxj/rpcx_consul/registry"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 0)
	fmt.Println(r.Body.Read(buf))
	fmt.Println(buf)

	fmt.Fprintln(w, "hello world")
}

func test() {
	a := "eyJzdHUiOjMsInNzc3MiOiJxaGhoaCIsInNza2siOjExMn0="
	bb, err := base64.StdEncoding.DecodeString(a)
	// b := bytes.Buffer{}
	// bb := []byte(a) s
	var entryInfos map[string]interface{}
	// c := json.NewDecoder(&b)
	// _,err := c.Buffered().Read(&a)
	json.Unmarshal([]byte(bb), &entryInfos)
	fmt.Println("...", err, entryInfos)

}
func main() {
	// test()
	// return
	demo.InitLog()
	demo.InitRegistry()
	// doServiceConvert()
	// conf := api.DefaultConfig()
	// conf.Address = "10.10.1.58:8500"
	// conf.Token = "1c67c216-a4ec-2235-da71-ad5d3f970280"
	// client, err := api.NewClient(conf)
	// if err != nil {
	// 	fmt.Println("client err: ", err)
	// }
	// fmt.Println(time.Now())
	// serv, err := client.Agent().ServicesWithFilter("dev in Tags")
	// if err != nil {
	// 	fmt.Println("client Service err: ", err)
	// }
	// fmt.Println("serv: ", serv)
	// for _, m := range serv {
	// 	if m.Service == "com.shining3d.app.auth.audit" {
	// 		fmt.Println("m: ", m.Service)
	// 		fmt.Println("m: ", m.Address)
	// 		fmt.Println("m: ", m.Port)
	// 	}
	// }
	// fmt.Println(time.Now())

	// var lastIndex uint64
	// ssss, smi, err := client.Catalog().Service("com.shining3d.app.auth.audit", "dev", &api.QueryOptions{
	// 	WaitIndex: lastIndex,
	// })
	// if err != nil {
	// 	fmt.Println("client Service err: ", err)
	// }
	// fmt.Println("ssss: ", ssss)
	// fmt.Println("smi: ", smi)
	// fmt.Println(time.Now())
	// rs, siii, err := client.Agent().AgentHealthServiceByName("com.shining3d.app.auth.audit")
	// if err != nil {
	// 	fmt.Println("client Service err: ", err)
	// }
	// fmt.Println("rs: ", rs)
	// fmt.Println("siii: ", siii)
	// for _, k := range siii {
	// 	if discovery.FindInStringSlice(k.Service.Tags, "dev") && k.AggregatedStatus == "passing" {
	// 		fmt.Println("k", k.Service.Service)
	// 		fmt.Println("k", k.Service.Address)
	// 		fmt.Println("k", k.Service.Port)
	// 	}
	// }
	// fmt.Println(time.Now())
}

func doServiceConvert() {
	fmt.Println("funcs: ", registry.GetServiceFunc(registry.ServiceFuncOBJ{
		ServiceFuncCommon: registry.ServiceFuncCommon{
			SFType: "func",
			SFName: "com.shining3d.app.list",
			SFCall: nil,
		},
		SFMeta: registry.ServiceFuncMeta{
			URLName:      "aaa",
			FuncName:     "bbb",
			URLPath:      "/ss",
			HTTPMethod:   "GET",
			AuthLevel:    "user",
			AuthPerms:    []string{"a.a"},
			ProductLines: []string{"dental"},
		},
	}))
	fmt.Println("class: ", registry.GetServiceFunc(registry.ServiceFuncOBJ{
		ServiceFuncCommon: registry.ServiceFuncCommon{
			SFType: "class",
			SFName: "com.shining3d.app.list",
			SFCall: nil,
		},
		SFMeta: registry.ServiceFuncMeta{
			Funcs: []string{"dental"},
		},
	}))
}
