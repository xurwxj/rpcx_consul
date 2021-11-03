package demo

import (
	"embed"

	"github.com/rs/zerolog"
	"github.com/xurwxj/gtils/net"
)

// Log global log instance
var Log *zerolog.Logger

const (
	timeDFormat = "2006-01-02T15:04:05.000Z"
	timeFormat  = "2006-01-02 15:04:05"
	timeFFormat = "2006-01-02T15:04:05Z"
	dateFormat  = "2006-01-02"
)

// CTLayout convert time layout
var CTLayout = []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05Z", "2006-01-02"}

var Templates embed.FS

// SMSTemplate used in aws
var SMSTemplate = map[string]string{}

var LocalesBytes []byte

// LocalesMap 国际化相关翻译字段，参考locales/locales.json这个文件
var LocalesMap map[string]map[string]string

var PushMqttManager *net.MqttManager
var PushTracingMqttManager *net.MqttManager
