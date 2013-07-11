package plugin

import (
	"irc"
	"irc/proto"
)

type adminStruct struct {
}

func init() {
	_, ok := PluginMap["admin"]
	if ok {
		panic("Shit happened!")
	}
	PluginMap["admin"] = adminParser
}

func adminParser(raw map[interface{}]interface{}) PluginInterface {
	var admin adminStruct
	return &admin
}

func (*adminStruct) Action(msg *proto.Message, conn irc.ConnInterface) {
	switch msg.Code {
	case "PRIVMSG":
		if msg.Arguments[0] == conn.GetCurrentNick() {
			conn.Privmsg("#ubuntu-cn", msg.Content)
		}
	default:
	}
}
