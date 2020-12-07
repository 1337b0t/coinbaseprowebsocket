package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

var server, serverAddr, port, productID string
var productIDList []string

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

	tcpAddr, _ := net.ResolveTCPAddr("tcp", serverAddr)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	conn.SetNoDelay(false)
	conn.SetWriteBuffer(10000)

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		//Message:
		//Market: BTC-USD
		//Price:  19000
		//Size:   1
		jsonMsg := fmt.Sprintf("[{market:%s,price:%s,size:%s}]", message.ProductID, message.Price, message.Size)
		conn.Write([]byte(jsonMsg))

		fmt.Println(jsonMsg)

	}

}

func main() {

	flag.StringVar(&server, "server", "127.0.0.1", "server address (local is default)")
	flag.StringVar(&port, "port", "8081", "the default port number is 8081")
	flag.StringVar(&productID, "products", "BTC-USD", "BTC-USD,ETH-USD")
	flag.Parse()

	//Server Adddress for TCP socket
	serverAddr = server + ":" + port

	//Create the product list for the Coinbase Pro Websocket
	split := strings.Split(productID, ",")
	for _, v := range split {
		productIDList = append(productIDList, v)
	}

	//Start Stream
	stream(productIDList)

}
