package main

import (
	"flag"
	"net"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

var (
	server, serverAddr, port, productID string
	productIDList                       []string
	conn                                net.Conn
)

func stream(productID []string, ch chan<- coinbasepro.Message) {

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
		panic("Script restart required")
	}

	go func() {

		keepAlive(wsConn, pongWait)
		for true {
			message := coinbasepro.Message{}
			if err := wsConn.ReadJSON(&message); err != nil {
				panic("Script restart required")
			}

			ch <- message

		}
	}()

}

func main() {

	flag.StringVar(&server, "server", "localhost", "server address (local is default)")
	flag.StringVar(&port, "port", "8081", "the default port number is 8081")
	flag.StringVar(&productID, "products", "BTC-USD,ETH-USD", "Example: BTC-USD,ETH-USD")
	flag.Parse()

	//Server Adddress for TCP socket
	serverAddr = server + ":" + port

	//Create the product list for the Coinbase Pro Websocket
	split := strings.Split(productID, ",")
	for _, v := range split {
		productIDList = append(productIDList, v)

	}

	//Setup channel for Coinbase Message
	ch := make(chan coinbasepro.Message)

	//Start Stream
	go stream(productIDList, ch)
	//Start TCP Server (testing)
	go tcpserver()
	//Start TCP Client
	time.Sleep(1 * time.Second)
	tcpclient(ch, serverAddr)

}
