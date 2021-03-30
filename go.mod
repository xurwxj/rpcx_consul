module github.com/xurwxj/rpcx_consul

go 1.15

require (
	github.com/hashicorp/consul/api v1.8.1
	github.com/smallnest/rpcx v0.0.0-20210302003640-3ac62d723635
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/xurwxj/gtils v1.0.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
