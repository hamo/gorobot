package irc

import (
	"time"
)

type ConnInterface interface {
	Connect() (err error)

	ActionLoop()

	PostConnect()
	Close() (err error)

	SetServerHost(host string)

	GetCurrentNick() string

	SetCurrentNick(nick string)

	SetLastAlive(time.Time)

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
