package irc

import (
	"crypto/tls"
	"net/textproto"
	//	"os"
	//	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
	"fmt"

	"conf"
	"irc/proto"
	"irc/action"
	"irc/plugin"
)

const (
	ConnAuthNone = 1 << iota
	ConnAuthPass
	ConnAuthNickServ
)

type Connection struct {
	conf *conf.ConnConf

	ServerHost string // server host name

	//	user      string
	//	host      string
	//	realname  string
	//	UseTLS    bool
	TLSConfig *tls.Config

	//	nick        string
	CurrentNick string

	//	Auth     int
	//	Password string

	//	ChanMap *ChannelMap

	Socket *textproto.Conn

	pwrite chan string
	pread  chan *proto.Message

	//	output, debug *log.Logger

	pingTicker *time.Ticker

	syncing    bool
	quiting    bool
	syncReader chan bool
	syncWriter chan bool
	syncPing   chan bool
	syncAction chan bool

	//	Error chan WorkerError
}

// type WorkerError struct {
// 	from string
// 	work string
// 	err  error
// 	//	arg string
// }

func NewConn(ca *conf.ConfStruct) (conn *Connection, err error) {
	conn = new(Connection)

	conn.conf = ca.Conn

	if conn.conf.SSL == true {
		conn.TLSConfig = new(tls.Config)
		conn.TLSConfig.InsecureSkipVerify = true
	}

	// if auth, err := ca.GetOptionOptional("CONNECTION", "auth", "none"); err != nil {
	// 	return nil, err
	// } else {
	// 	for _, s := range strings.Split(auth, ",") {
	// 		s = strings.TrimSpace(strings.ToLower(s))
	// 		switch s {
	// 		case "none":
	// 			conn.Auth |= ConnAuthNone
	// 		case "pass":
	// 			conn.Auth |= ConnAuthPass
	// 		case "nickserv":
	// 			conn.Auth |= ConnAuthNickServ
	// 		default:
	// 		}
	// 	}
	// }

	// if pass, err := ca.GetOptionOptional("CONNECTION", "pass", ""); err != nil {
	// 	return nil, err
	// } else {
	// 	conn.Password = pass
	// }

	// if nick, err := ca.GetOption("CONNECTION", "nick"); err != nil {
	// 	return nil, err
	// } else {
	// 	conn.nick = nick
	// }

	// if output, err := ca.GetOptionOptional("CONNECTION", "output", ""); err != nil {
	// 	return nil, err
	// } else {
	// 	prefix := strings.Join([]string{"[", conn.addr, "INFO", "]"}, " ")
	// 	conn.output, err = log.CreateLoggerFromFileString(output, prefix, log.LstdFlags)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// if debug, err := ca.GetOptionOptional("CONNECTION", "debug", ""); err != nil {
	// 	return nil, err
	// } else {
	// 	prefix := strings.Join([]string{"[", conn.addr, "DEBUG", "]"}, " ")
	// 	conn.debug, err = log.CreateLoggerFromFileString(debug, prefix, log.LstdFlags)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// cm, err := NewChannelMap(ca)
	// if err != nil {
	// 	return nil, err
	// }
	// conn.ChanMap = cm
	return conn, nil
}

func (conn *Connection) Connect() (err error) {
	addr := strings.Join([]string{conn.conf.Server, strconv.Itoa(conn.conf.Port)}, ":")
	if conn.conf.SSL {
		socket, err := tls.Dial("tcp", addr, conn.TLSConfig)
		if err != nil {
			return err
			//FIXME
		}
		conn.Socket = textproto.NewConn(socket)
	} else {
		conn.Socket, err = textproto.Dial("tcp", addr)
		if err != nil {
			return err
			//FIXME
		}
	}

	return nil
}

func (conn *Connection) initConn() {
	conn.CurrentNick = conn.conf.Nick
	
	conn.pwrite = make(chan string, 64)
	conn.pread = make(chan *proto.Message, 64)

	conn.syncReader = make(chan bool)
	conn.syncWriter = make(chan bool)
	conn.syncPing = make(chan bool)
	conn.syncAction = make(chan bool)

	// conn.Error = make(chan WorkerError)

	conn.pingTicker = time.NewTicker(15 * time.Minute)
}

func (conn *Connection) closeConn() {
	close(conn.pwrite)
	close(conn.pread)

	close(conn.syncReader)
	close(conn.syncWriter)
	close(conn.syncPing)
	close(conn.syncAction)

	// close(conn.Error)

	conn.pingTicker.Stop()
}

func (conn *Connection) PostConnect() {
	conn.initConn()

	go conn.WriteLoop()
	go conn.ReadLoop()
	go conn.PingLoop()
	go conn.ActionLoop()

	conn.User(conn.conf.User, "test", "0.0.0.0", conn.conf.Realname)
	conn.Nick(conn.conf.Nick)

	//	if (conn.Auth&ConnAuthPass) != 0 && conn.Password != "" {
	//		conn.Pass(conn.Password)
	//	}

	// conn.ChanMap.RLock()
	// for k, _ := range conn.ChanMap.CM {
	// 	conn.Join(k)
	// }
	// conn.ChanMap.RUnlock()
	conn.Join("#ubuntu-cn")
}

func (conn *Connection) Close() (err error) {
	return nil
}

func (conn *Connection) SetServerHost(host string) {
	conn.ServerHost = host
}

func (conn *Connection) GetCurrentNick() string {
	return conn.CurrentNick
}

func (conn *Connection) SetCurrentNick(nick string) {
	conn.CurrentNick = nick
}

func (conn *Connection) ReadLoop() {
	for {
		select {
		// case sync := <-conn.syncReader:
		// 	if sync {
		// 		for conn.syncing {
		// 			runtime.Gosched()
		// 		}
		// 		if conn.quiting {
		// 			return
		// 		}
		// 	}
		default:
			if s, err := conn.Socket.Reader.ReadLine(); err != nil {
				//				conn.Error <- WorkerError{"Reader", "Read Socket", err}
			} else {
				if msg, err := proto.Parse(s); err != nil {
					// 	conn.Error <- WorkerError{"Reader", "Parse message", err}
				} else {
					conn.pread <- msg
				}
			}
		}
	}
}

func (conn *Connection) PingLoop() {
	for {
		select {
		case <-conn.pingTicker.C:
			conn.Ping(string(time.Now().UnixNano()))
			//Try to recapture nickname if it's not as configured.
			if conn.conf.Nick != conn.CurrentNick {
				conn.Nick(conn.conf.Nick)
			}
		case sync := <-conn.syncPing:
			if sync {
				for conn.syncing {
					runtime.Gosched()
				}
				if conn.quiting {
					return
				}
			}
		default:
			runtime.Gosched()
		}
	}
}

func (conn *Connection) WriteLoop() {
	for {
		select {
		case msg := <-conn.pwrite:
			conn.Socket.PrintfLine("%s", msg)
		case sync := <-conn.syncWriter:
			if sync {
				for conn.syncing {
					runtime.Gosched()
				}
				if conn.quiting {
					return
				}
			}
		default:
			runtime.Gosched()
		}
	}
}

func (conn *Connection) ActionLoop() {
	for {
		select {
		case sync := <-conn.syncAction:
			if sync {
				for conn.syncing {
					runtime.Gosched()
				}
				if conn.quiting {
					return
				}
			}
		case msg := <-conn.pread:
			go action.Action(msg, conn)
			go plugin.Action(msg, conn)
		default:
			runtime.Gosched()
		}
	}
}


func (conn *Connection) Join(channel string) {
	conn.pwrite <- fmt.Sprintf("JOIN %s", channel)
}

func (conn *Connection) Quit() {
	conn.pwrite <- fmt.Sprintf("QUIT")
}

func (conn *Connection) Part(channel string) {
	conn.pwrite <- fmt.Sprintf("PART %s", channel)
}

func (conn *Connection) Ping(time string) {
	conn.pwrite <- fmt.Sprintf("PING %s", time)
}

func (conn *Connection) Pong(time string) {
	conn.pwrite <- fmt.Sprintf("PONG %s", time)
}

func (conn *Connection) Notice(target, message string) {
	conn.pwrite <- fmt.Sprintf("NOTICE %s :%s", target, message)
}

func (conn *Connection) Privmsg(target, message string) {
	conn.pwrite <- fmt.Sprintf("PRIVMSG %s :%s", target, message)
}

func (conn *Connection) User(user, host, server, realname string) {
	conn.pwrite <- fmt.Sprintf("USER %s %s %s :%s", user, host, server, realname)
}

func (conn *Connection) Pass(pass string) {
	conn.pwrite <- fmt.Sprintf("PASS %s", pass)
}

func (conn *Connection) Nick(n string) {
	conn.pwrite <- fmt.Sprintf("NICK %s", n)
}

func (conn *Connection) SendRawf(format string, a ...interface{}) {
	conn.pwrite <- fmt.Sprintf(format, a...)
}
