//http://www.networksorcery.com/enp/protocol/irc.htm

package proto

import (
	"strings"
)

type Message struct {
	Raw string

	Source string //<nick>[!user][@host]
	Nick   string
	User   string
	Host   string

	Code      string
	Content   string
	Arguments []string
}

func Parse(line string) (msg *Message, e error) {
	msg = new(Message)
	msg.Raw = line

	if line == "" {
		return
	}
	if strings.HasPrefix(line, ":") {
		if i := strings.Index(line, " "); i != -1 {
			// Server messages
			msg.Source = line[1:i]
			line = line[i+1:]
		} else {
			// Misformat messages
		}

		if i, j := strings.Index(msg.Source, "!"), strings.Index(msg.Source, "@"); i > -1 && j > -1 {
			msg.Nick = msg.Source[0:i]
			msg.User = msg.Source[i+1 : j]
			msg.Host = msg.Source[j+1:]
		} else if i > -1 && j == -1 {
			msg.Nick = msg.Source[0:i]
			msg.User = msg.Source[i+1:]
		} else if i == -1 && j > -1 {
			msg.Nick = msg.Source[0:j]
			msg.Host = msg.Source[j+1:]
		} else {
			msg.Nick = msg.Source
		}
	}

	args := strings.SplitN(line, " :", 2)
	if len(args) > 1 {
		msg.Content = args[1]
	}

	args = strings.Split(args[0], " ")
	msg.Code = strings.ToUpper(args[0])

	if len(args) > 1 {
		msg.Arguments = args[1:]
	}
	/* XXX: len(args) == 0: args should be empty */
	return msg, nil
}
