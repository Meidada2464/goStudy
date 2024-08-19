package main

import "github.com/lesismal/nbio"

func main() {
	tcpServerStart()
}

func tcpServerStart() {
	// config tcp servers
	nbio.NewEngine(nbio.Config{})
}
