package plugin

import (
	"irc"
	"irc/proto"

	"utils"
)

type adminStruct struct {
	adminNicks []string
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

	if v, ok := raw["admins"]; ok {
		raw := v.([]interface{})
		for _, nick := range raw {
			admin.adminNicks = append(admin.adminNicks, nick.(string))
		}
	}
	
	return &admin
}

func (as *adminStruct) Action(msg *proto.Message, conn irc.ConnInterface) {
	switch msg.Code {
	case "PRIVMSG":
		if msg.Arguments[0] == conn.GetCurrentNick() {
			if utils.StringInSlice(msg.Nick, as.adminNicks) {
				conn.Privmsg("#gorobot", msg.Content)
			} else {
				conn.Privmsg(msg.Nick, "你不是我主人,你是坏人!")
			}
		}
	default:
	}
}
