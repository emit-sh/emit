package main

import (
	"github.com/emit-sh/emit/server"
	"fmt"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		fmt.Print(err)
		return
	}

	server.Start()
}
