package conf

import (
	"errors"
	"io/ioutil"
	"launchpad.net/goyaml"
	"path/filepath"

	"utils"
)

const (
	defaultNick = "gorobot"
	defaultUser = "gorobot"
)

type ConfStruct struct {
	Filename string
	Conn     *ConnConf
	Channels []channelConf
	Plugins  []map[interface{}]interface{}
}

type ConnConf struct {
	Server   string
	Port     int
	SSL      bool
	Nick     string
	User     string
	Realname string
}

type channelConf struct {
	Name string
}

func ParseFile(filename string) (*ConfStruct, error) {
	var err error
	var content []byte
	var raw map[interface{}]interface{}
	var cs *ConfStruct

	utils.BUG_ON(!filepath.IsAbs(filename), "filename is not abs path")

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		goto e
	}

	err = goyaml.Unmarshal(content, &raw)
	if err != nil {
		goto e
	}

	cs, err = parseContent(raw)
	if err != nil {
		goto e
	}

	cs.Filename = filename
	return cs, nil

e:
	return nil, err
}

func parseContent(raw map[interface{}]interface{}) (*ConfStruct, error) {
	cs := new(ConfStruct)
	for k, v := range raw {
		switch k {
		case "conn":
			r := v.(map[interface{}]interface{})
			cc, err := parseConn(r)
			if err != nil {
				return nil, err
			}
			cs.Conn = cc
		case "channels":
			r := v.([]interface{})
			cc, err := parseChannels(r, "")
			if err != nil {
				return nil, err
			}
			cs.Channels = cc
		case "plugins":
			r := v.([]interface{})
			cs.Plugins = parsePlugins(r)
		}
	}
	return cs, nil
}

func parseConn(raw map[interface{}]interface{}) (*ConnConf, error) {
	cc := new(ConnConf)
	for k, v := range raw {
		switch k {
		case "server":
			cc.Server = v.(string)
		case "port":
			cc.Port = v.(int)
		case "ssl":
			cc.SSL = v.(bool)
		case "nick":
			cc.Nick = v.(string)
		case "user":
			cc.User = v.(string)
		case "realname":
			cc.Realname = v.(string)
		default:
			//LOG
			continue
		}
	}

	//sanity check
	if cc.Server == "" {
		return nil, errors.New("No server address found")
	}
	if cc.Port == 0 {
		cc.Port = 6667
	}
	if cc.Nick == "" {
		cc.Nick = defaultNick
	}
	if cc.User == "" {
		cc.User = defaultUser
	}
	if cc.Realname == "" {
		cc.Realname = cc.User
	}

	return cc, nil
}

func parseChannels(raw []interface{}, name string) ([]channelConf, error) {
	ccs := make([]channelConf, 0, 5)
	for _, v := range raw {
		value := v.(map[interface{}]interface{})
		if name != "" && name != value["name"] {
			continue
		}
		cc, _ := parseOneChannel(value)
		ccs = append(ccs, cc)

	}

	utils.BUG_ON((name != "" && len(ccs) != 1), "name special channel parse not return 1 channel")

	return ccs, nil
}

func parseOneChannel(raw map[interface{}]interface{}) (channelConf, error) {
	var cc channelConf
	for k, v := range raw {
		switch k {
		case "name":
			cc.Name = "#" + v.(string)
		}
	}
	return cc, nil
}

func parsePlugins(raw []interface{}) []map[interface{}]interface{} {
	var r []map[interface{}]interface{}
	for _, v := range raw {
		r = append(r, v.(map[interface{}]interface{}))
	}
	return r
}
