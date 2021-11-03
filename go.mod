module github.com/xurwxj/rpcx_consul

go 1.16

require (
	git.shining3d.com/cloud/util v0.0.0-00010101000000-000000000000
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/hashicorp/consul/api v1.11.0
	github.com/rs/zerolog v1.21.0
	github.com/smallnest/rpcx v1.6.11
	github.com/soheilhy/cmux v0.1.4
	github.com/xurwxj/gtils v1.0.7
	github.com/xurwxj/viper v1.7.1

)

replace (
	git.shining3d.com/cloud/util => ../util
	git.shining3d.com/cloud/util/service => ../util/service
)
