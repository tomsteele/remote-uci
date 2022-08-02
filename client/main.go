package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Addr string
	Key  string
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("remote_uci.toml", &conf)
	if err != nil {
		log.Println("error reading toml")
		log.Println(err)
		os.Exit(1)
	}
	conn, err := net.Dial("tcp", conf.Addr)
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
