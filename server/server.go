package server

import (
	"bufio"
	"fmt"
	"github.com/ryym/vimdbb/vimch"
	"net"
	"os"
)

func Listen(ln net.Listener) error {
	sysChan := make(chan string)
	connChan := make(chan net.Conn)

	sysConn, err := ln.Accept()
	if err != nil {
		return err
	}

	go listenSysConn(sysChan, sysConn)
	go listenUserConns(connChan, ln)
	serve(sysChan, connChan)
	return nil
}

func listenSysConn(ch chan string, conn net.Conn) {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		mes := sc.Text()
		m := vimch.DecodeMessage(mes)
		ch <- m.Command
	}
}

func listenUserConns(ch chan net.Conn, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			ch <- conn
		}
	}
}

func serve(sysChan chan string, connChan chan net.Conn) {
	for {
		select {
		case command := <-sysChan:
			if command == "KILL" {
				return
			}
		case conn := <-connChan:
			defer conn.Close()
			go handleUserConn(conn)
		}
	}
}

func handleUserConn(conn net.Conn) {
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		message := sc.Text()
		fmt.Println(message)

		result, err := handleUserMessage(message)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if result != nil {
			conn.Write(result)
		}
	}

	fmt.Println("disconnected")
}
