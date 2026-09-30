package main

import (
	"fmt"
	"net"
)

type CmdHandler func(args []string, conn net.Conn) error

var handlers = map[string]CmdHandler{
	"ping": ping,
}

func ping(args []string, conn net.Conn) error {
	resWithStr(conn, "ping ok")
	return nil
}

func resWithStr(conn net.Conn, s string) {
	fmt.Fprintf(conn, "+%s\r\n", s)
}
