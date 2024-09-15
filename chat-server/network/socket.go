package network

import (
	"chat_server_golang/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

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

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

// 클라이언트가 들어오는 메시지를 읽는 함수
func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		var msg *Message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			panic(err)
		} else {
			log.Println("READ : ", msg, "client", c.Name)
			log.Println()
			msg.Time = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

// 클라이언트가 메시지를 전송하는 함수
func (c *Client) Write() {
	defer c.Socket.Close()

	for msg := range c.Send {
		log.Println("WRITE : ", msg, "client", c.Name)
		log.Println()
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) RunInit() {
	// Room에 있는 모든 채널값들을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Socket: socket,
		Send:   make(chan *Message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client

	defer func() { r.Leave <- client }()

	go client.Write()

	client.Read()
}
