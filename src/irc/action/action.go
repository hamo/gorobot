package action

import (
	"irc/proto"
	"irc"
	"fmt"
)

func Action(msg *proto.Message, conn irc.ConnInterface) {
	switch msg.Code {
	case "001":
		conn.SetServerHost(msg.Source)
	case "NICK":
		if msg.Nick == conn.GetCurrentNick() {
			conn.SetCurrentNick(msg.Content)
		}
	case "PING":
		conn.Pong(msg.Content)
	// case "PRIVMSG":
	// 	for _, v := range msg.Arguments {
	// 		if strings.HasPrefix(v, "#") {
	// 			conn.ChanMap.RLock()
	// 			if v, ok := conn.ChanMap.CM[v]; ok {
	// 				v.output.Printf("<%s> %s\n", msg.Nick, msg.Content)
	// 			}
	// 			conn.ChanMap.RUnlock()
	// 		} else {
	// 			conn.output.Printf("PRIVMSG <%s> %s\n", msg.Nick, msg.Content)
	// 		}
	// 	}
	case "PONG":
		conn.SetServerHost(msg.Arguments[0])
	case "NOTICE":
		fallthrough
	default:
		fmt.Printf("%+v\n", msg)
	}
}
