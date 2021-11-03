package demo

import (
	"context"

	"git.shining3d.com/cloud/util/service"
	"github.com/xurwxj/rpcx_consul/registry"
)

var serviceFuncs = []registry.ServiceFuncItem{
	registry.GetServiceFunc(registry.ServiceFuncOBJ{
		ServiceFuncCommon: registry.ServiceFuncCommon{SFType: "func", SFName: "com.shining3d.sm.log", SFCall: softModularUseLog},
		SFMeta:            registry.ServiceFuncMeta{URLName: "softModularUseLog", FuncName: "SoftModularUseLog", URLPath: "/sm/log", HTTPMethod: "get", AuthLevel: "userPerm", AuthPerms: []string{"dental.view_sm"}, ProductLines: []string{"dental"}},
	}),
	// registry.GetServiceFunc(registry.ServiceFuncOBJ{
	// 	ServiceFuncCommon: registry.ServiceFuncCommon{SFType: "func", SFName: "com.shining3d.sm.use", SFCall: SoftModularUse},
	// 	SFMeta:            registry.ServiceFuncMeta{URLName: "softModularUse", FuncName: "SoftModularUse", URLPath: "/sm/use", HTTPMethod: "post", AuthLevel: "user", ProductLines: []string{"dental"}},
	// }),
}

func SoftModularUse(ctx context.Context, args *service.Args, reply *service.Replies) error {
	//1参数校验 不符合的设置成默认数据 {
	return nil
}

func softModularUseLog(ctx context.Context, args *service.Args, reply *service.Replies) error {
	//1参数校验 不符合的设置成默认数据 {
	return nil
}
