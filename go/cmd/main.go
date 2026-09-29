package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		slog.Error("failed to start tcp listener", "error", err.Error())
		return
	}
	slog.Info("listening for TCP on :6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("failed to start accepting tcp conns", "error", err.Error())
			return
		}
		go handleConn(conn)
	}
}

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// *2\r\n$3\r\nGET\r\n$3\r\nkey\r\n

// *2\r\n
// $3\r\n
// GET\r\n
// $3\r\n
// key\r\n

// = ["GET","key"].

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		res, err := parse(reader, conn)
		if err != nil {
			slog.Error("err", "err", err)
			conn.Write([]byte("+err\r\n"))
		}
		fmt.Printf("val: |%v|\n", res)
		conn.Write([]byte("+OK\r\n"))

	}
}

func parse(reader *bufio.Reader, conn net.Conn) ([]string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch rune(b) {
	case ARRAY:
		// parse every item of arr, if they bulks the switch statement will handle them
		// and we accumulate results in arr
		l := readInt(reader)
		cmds := make([]string, 0, l)
		for range l {
			cmdParts, err := parse(reader, conn)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, cmdParts...)
		}
		return cmds, nil
	case BULK:
		l := readInt(reader)
		buf := make([]byte, l)
		// read command str until \r\n
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}
		// discard the \r\n left
		reader.Discard(2)
		return []string{string(buf)}, nil
	}
	return nil, nil
}

func readInt(reader *bufio.Reader) int64 {
	i, _ := reader.ReadString('\n')
	i = i[:len(i)-2]
	l, _ := strconv.ParseInt(string(i), 0, 64)
	if l <= 0 {
		return 0
	}
	return l
}
