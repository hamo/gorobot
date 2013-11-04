package main

import (
	"conf"
	"flag"
	"path/filepath"
	"runtime"
	"strings"
	"utils"

	"irc/conn"
)

var (
	confDir   = flag.String("config", "conf", "config dir")
	debug     = flag.Bool("debug", true, "debug switch")

//	FIXME: We need many co-run goroutines to make real-time reaction, set default GOMAXPROCS to 64
//	threadNum = flag.Int("threadNum", runtime.NumCPU()+1, "thread number")
	threadNum = flag.Int("threadNum", 64, "thread number")
)

var (
	connSlice []*conn.Connection
)

func main() {

	runtime.GOMAXPROCS(*threadNum)
	flag.Parse()

	confs, err := utils.AllFilesUnderDir(*confDir)
	if err != nil {
		panic(err)
	}

	if len(confs) == 0 {
		panic("No config file found.")
	}

	for _, v := range confs {
		v = filepath.Join(*confDir, v)
		if !strings.HasSuffix(v, ".yaml") {
			continue
		}

		if !filepath.IsAbs(v) {
			v, err = filepath.Abs(v)
		}
		if err != nil {
			panic(err)
		}

		cs, err := conf.ParseFile(v)
		if err != nil {
			continue
		}

		oneConn, err := conn.NewConn(cs)
		if err != nil {
			continue
		}

		connSlice = append(connSlice, oneConn)
	}

	if len(connSlice) == 0 {
		panic("No vaild conn found.")
	}

	for _, c := range connSlice {
		println("Start connect to server")
		err = c.Connect()
		if err != nil {
			continue
		}

		println("Post connect to server")
		c.PostConnect()
	}

	for {
		for _, c := range connSlice {
			select {
			case err := <-c.Error:
				panic(err.Err)
			default:
				runtime.Gosched()
			}
		}
	}
}
