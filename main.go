package main

import (
	"log"

	"chat_server_golang/network"
)

func init() {
	log.Println("먼저 시작됩니다.")

	network := network.NewServer()
	network.StartServer()
}

func main() {
	log.Println("----")
}
