package main

import (
	"fmt"
	"github.com/ryym/vimdbb/server"
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
	server.Listen(ln)
	fmt.Println("server stopped")
}
