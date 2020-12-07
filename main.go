package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

var (
	servAddr string
)

func stream(productID []string) {

	var wsDialer websocket.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{
				Name:       "heartbeat",
				ProductIds: productID,
			},
			coinbasepro.MessageChannel{
				Name:       "ticker",
				ProductIds: productID,
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	conn.SetNoDelay(false)
	conn.SetWriteBuffer(10000)

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		jsonMsg := fmt.Sprintf("[{market:%s,price:%s,quantity:%s", message.ProductID, message.Price, message.Size)
		start := time.Now()
		conn.Write([]byte(jsonMsg))
		fmt.Println("took:", time.Since(start))

	}

}

func main() {
	flag.StringVar(&servAddr, "server", "127.0.0.1", "server address (local is default)")

	flag.Parse()

}
