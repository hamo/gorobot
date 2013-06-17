package irc

import (
)

type ConnInterface interface {
	Connect() (err error)

//	initConn()
//	closeConn()

	ActionLoop()

	PostConnect()
	Close() (err error)

	SetServerHost(host string)

	GetCurrentNick() string

	SetCurrentNick(nick string)

	//proto
	Join(channel string)
	Quit()
	Part(channel string)
	Ping(time string)
	Pong(time string)
	Notice(target, message string)
	Privmsg(target, message string)
	User(user, host, server, realname string)
	Pass(pass string)
	Nick(n string)
}
