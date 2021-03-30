module github.com/xurwxj/rpcx_consul

go 1.16

require (
	github.com/hashicorp/consul/api v1.8.1
	github.com/smallnest/rpcx v0.0.0-20210329112732-c584448849f9
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/xurwxj/gtils v1.0.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
