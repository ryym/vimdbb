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

		action := vimch.DecodeMessage(message)

		result := run(action.Payload.Query)
		resJson, err := vimch.EncodeMessage(action.Id, vimdbb.Result{
			Rows: result,
		})

		if err != nil {
			panic(err.Error())
		}
		conn.Write(resJson)
	}

	fmt.Println("END!")
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
