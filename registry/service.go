package registry

import "github.com/xurwxj/gtils/base"

// GetServiceFunc convert service obj to service definition
func GetServiceFunc(s ServiceFuncOBJ) (sf ServiceFuncItem) {
	sf.ServiceFuncCommon = s.ServiceFuncCommon
	meta := base.GetStringFromInterface(s.SFMeta)
	switch s.SFType {
	case "func":
		sf.SFMeta = `httpInfo=` + meta
	case "class":
		sf.SFMeta = `funcs=` + meta
	}
	return
}

// ServiceFuncItem service or func obj struct, used for rpcx in registry
type ServiceFuncItem struct {
	ServiceFuncCommon
	// 压缩后的meta信息
	// funcs=["APICheck","APIUserCheck","APIUserPermCheck"]
	// httpInfo={"name":"resourceKeyList","funcName":"ResourceKeyList","path":"/settings/r/:nameCode/:keyName","method":"GET","auth":"api","productLines":["dentalscan","scan","dlp","thirdpartner"]}
	SFMeta string
}

// ServiceFuncCommon used in service definition
type ServiceFuncCommon struct {
	// 	SFType 有两种值：
	//   - func 表示单个函数作为服务，适合某个接口就是单独的http服务
	//   - class 适合把一些函数做集合，挂到统一的struct下，这些函数一般是不提供http服务的，只用于服务间的调用
	SFType string
	// 服务名，一般是用java类名的定义方式，比如com.shining3d.app.app
	SFName string
	// 真正的服务执行方法或类
	SFCall interface{}
}

// ServiceFuncOBJ used in service definition
type ServiceFuncOBJ struct {
	ServiceFuncCommon
	SFMeta ServiceFuncMeta
}

// ServiceFuncMeta used in service definition
type ServiceFuncMeta struct {
	// name: 服务或接口唯一名，用英文，一般用于外部接入时避免网址写死，会在cdn的json作为key
	URLName string `json:"name,omitempty"`
	// funcName: 方法函数名
	FuncName string `json:"funcName"`
	// http路径，注意如果重复，则重复的这些服务只有1个有效，无序
	URLPath    string `json:"path,omitempty"`
	HTTPMethod string `json:"method,omitempty"`
	// 	auth: 指服务的验证方式
	//   - 为空表示不用验证
	//   - api 表示验证appID
	//   - user 表示验证用户有效性
	//   - userPerm 表示验证用户及权限是否有效
	//   - apiToken 表示验证appID及token是否有效
	AuthLevel string `json:"auth,omitempty"`
	// 权限字符串数组，一般只要拥有其中1项即可
	AuthPerms []string `json:"perms,omitempty"`
	// 产品线字符串数组，设置后，对应cdn的产品线json网址中会出现定义的服务名和网址
	ProductLines []string `json:"productLines,omitempty"`
	Funcs        []string `json:"funcs,omitempty"`
}
