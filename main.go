package main

import "emit/server"

func main() {
	server, err := server.NewServer()
	if err != nil {
		return
	}

	server.Start()
}
