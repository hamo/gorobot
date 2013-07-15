package plugin

import (
	"irc"
	"irc/proto"
	"utils"
)

type adminStruct struct {
	passwd        string
	adminVerified map[string]bool
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
	admin.adminVerified = make(map[string]bool)

	if v, ok := raw["passwd"]; ok {
		admin.passwd = v.(string)
	}

	if v, ok := raw["admins"]; ok {
		raw := v.([]interface{})
		for _, nick := range raw {
			admin.adminVerified[nick.(string)] = false
		}
	}

	return &admin
}

func (as *adminStruct) Action(msg *proto.Message, conn irc.ConnInterface) {
	switch msg.Code {
	case "PRIVMSG":
		if msg.Arguments[0] == conn.GetCurrentNick() { // private message
			if v, ok := as.adminVerified[msg.Nick]; ok {
				command := utils.SplitAndTrimN(msg.Content, " ", 2)
				if v {
					switch command[0] {
					case "privmsg":
						args := utils.SplitAndTrimN(command[1], " ", 2)
						channel := args[0]
						msg := args[1]
						conn.Privmsg(channel, msg)
					default:
					}
				} else {
					// only allow VERIFY message
					if command[0] == "verify" && command[1] == as.passwd {
						as.adminVerified[msg.Nick] = true
						conn.Privmsg(msg.Nick, "VERIFIED")
					}
				}
			} else {
				// Do nothing if not an administrator
			}
		}
	case "QUIT":
		if _, ok := as.adminVerified[msg.Nick]; ok {
			as.adminVerified[msg.Nick] = false
		}
	default:
	}
}
