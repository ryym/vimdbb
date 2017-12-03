package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ryym/vimdbb/server"
)

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer ln.Close()

	fmt.Println("server started")
	server.Listen(ln)
	fmt.Println("server stopped")
}
