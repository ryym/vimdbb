package main

import (
	"bufio"
	"fmt"
	"github.com/ryym/vimdbb"
	"github.com/ryym/vimdbb/mysql"
	"github.com/ryym/vimdbb/vimch"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer ln.Close()

	fmt.Println("server started")

	sysChan := make(chan string)
	userChan := make(chan net.Conn)

	sysConn, err := ln.Accept()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("sys connection start")
	go listenSysCommands(sysChan, sysConn)
	go listenUserConns(userChan, ln)
	serve(sysChan, userChan)

	fmt.Println("server stopped")
}

func serve(sysChan chan string, userChan chan net.Conn) {
	conns := make([]net.Conn, 0)
	for {
		select {
		case command := <-sysChan:
			keep := handleSysCommand(command, conns)
			if !keep {
				return
			}

		case conn := <-userChan:
			defer conn.Close()
			conns = append(conns, conn)
			go handleUserConn(conn)

		default:
		}
	}
}

func listenSysCommands(ch chan string, conn net.Conn) {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		message := sc.Text()
		_, command, _ := vimch.DecodeMessage(message)
		fmt.Println(message, command)
		ch <- command
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

func handleSysCommand(command string, conns []net.Conn) bool {
	switch command {
	case "KILL":
		for _, conn := range conns {
			conn.Close()
		}
		return false
	}
	return true
}

func handleUserConn(conn net.Conn) {
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		message := sc.Text()
		fmt.Println(message)

		result, err := handleMessage(message)
		if err != nil {
			panic(err.Error())
		}
		conn.Write(result)
	}

	fmt.Println("disconnected")
}

func handleMessage(message string) ([]byte, error) {
	id, action, payload := vimch.DecodeMessage(message)
	switch action {
	case "Query":
		queryP := vimdbb.QueryPayload{}
		vimch.DecodePayload(payload, &queryP)
		return handleQuery(id, queryP)
	}
	panic("Unknown action " + action)
}

func handleQuery(id float64, p vimdbb.QueryPayload) ([]byte, error) {
	result := run(p.Query)
	return vimch.EncodeMessage(id, vimdbb.Result{
		Rows: result,
	})
}

func run(queryStr string) string {
	db, err := mysql.Open("root:root@/sample")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query(queryStr)
	if err != nil {
		panic(err.Error())
	}

	return mysql.ResultToString(result)
}
