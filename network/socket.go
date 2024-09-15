package network

import (
	"github.com/gorilla/websocket"
)

type Room struct {
	Forward chan *Message // 수신되는 메시지를 보관하는 값
	// 들어오는 메시지를 다른 클라이언트들에게 전송을 합니다.

	Join  chan *Client // Socket이 연결되는 경우에 작용
	Leave chan *Client // Scoket이 끊어지는 경우에 대해서 작용

	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

type Message struct {
	Name    string
	Message string
	Time    int64
}

type Client struct {
	Send   chan *Message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}
