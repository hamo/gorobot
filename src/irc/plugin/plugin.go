package plugin

import (
	"irc"
	"irc/proto"
)

var PluginMap map[string]PluginParser = make(map[string]PluginParser)

type PluginInterface interface {
	Action(*proto.Message, irc.ConnInterface)
}

type PluginParser func(map[interface{}]interface{}) PluginInterface
