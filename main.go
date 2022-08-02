package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func croak(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("remote-uci ip:port engine_binary")
		os.Exit(1)
	}
	addr := os.Args[1]
	engine := os.Args[2]
	log.Printf("Listening on %s\n", addr)
	log.Printf("Using engine %s\n", engine)
	l, err := net.Listen("tcp", addr)
	croak(err)
	for {
		conn, err := l.Accept()
		croak(err)
		log.Printf("New connection from %s\n", conn.RemoteAddr())
		p := exec.Command(engine)
		stdin, err := p.StdinPipe()
		croak(err)
		stdout, err := p.StdoutPipe()
		croak(err)
		croak(p.Start())

		go func() {
			for {
				buff := make([]byte, 2048)
				n, err := conn.Read(buff)
				if err != nil {
					log.Println("There was an error in reading from the connection")
					conn.Close()
					p.Process.Kill()
					break
				}
				stdin.Write(buff[:n])
			}
		}()

		go io.Copy(conn, stdout)
		if err := p.Wait(); err != nil {
			log.Printf("Process was killed with error: %s\n", err.Error())
		}
	}
}
