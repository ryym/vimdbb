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

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		defer conn.Close()

		fmt.Println("accepted")
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		fmt.Println("scanned")
		message := sc.Text()

		fmt.Println(message)

		result, err := handleMessage(message)
		if err != nil {
			panic(err.Error())
		}
		conn.Write(result)
	}

	fmt.Println("END!")
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
