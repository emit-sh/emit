package main

import (
	"emit/server"
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
