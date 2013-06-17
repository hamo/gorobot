package plugin

import (
	"irc"
	"irc/proto"
	"strings"
)


func Action(msg *proto.Message, conn irc.ConnInterface) {
	if strings.Contains(msg.Content, "test") {
		conn.Privmsg("#ubuntu-cn", "test failed.")
	}
}
	











