package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ryym/vimdbb"
	"github.com/ryym/vimdbb/formatter"
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
		rawMessage := sc.Text()
		m := vimch.DecodeMessage(rawMessage)
		fmt.Println(rawMessage)
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

func handleUserMessage(rawMessage string) ([]byte, error) {
	m := vimch.DecodeMessage(rawMessage)
	result, err := handleUserCommand(m)

	if err != nil {
		res := make(map[string]interface{})
		res["Command"] = "ERR"
		res["Result"] = err.Error()
		encm, encErr := vimch.EncodeMessage(m.ID, res)
		if encErr != nil {
			return nil, encErr
		}
		return encm, err
	}

	if result != nil {
		res := make(map[string]interface{})
		res["Command"] = m.Command
		res["Result"] = result
		return vimch.EncodeMessage(m.ID, res)
	}

	return nil, nil
}

func handleUserCommand(m *vimdbb.Message) (interface{}, error) {
	switch m.Command {
	case "Query":
		queryP := vimdbb.QueryPayload{}
		vimch.DecodePayload(m.Payload, &queryP)
		return handleQuery(queryP)
	}

	return nil, errors.New("Unknown command " + m.Command)
}

func handleQuery(p vimdbb.QueryPayload) (*vimdbb.Result, error) {
	db, err := mysql.Open("root:root@/sample")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	queryResult, err := db.Query(p.Query)
	if err != nil {
		return nil, err
	}

	rows := formatter.ResultToString(queryResult)
	result := vimdbb.Result{QueryID: p.QueryID, Rows: rows}
	return &result, nil
}
