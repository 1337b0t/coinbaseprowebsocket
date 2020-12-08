package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

var server, port, productID string
var serverAddr string
var productIDList []string

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
		println(err.Error())
	}

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		ch <- message

	}

}

func sender(ch <-chan coinbasepro.Message) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)

	if err != nil {
		println("ResolveTCPAddr Failed:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		println("DialTCP Failed:", err.Error())
	}

	for {

		select {
		case v, err := <-ch:
			if !err {
				return
			}
			//Message:
			//Market: BTC-USD
			//Price:  19000
			//Size:   1
			jsonMsg := fmt.Sprintf("[{market:%s,price:%s,size:%s}]", v.ProductID, v.Price, v.LastSize)
			conn.Write([]byte(jsonMsg))
			fmt.Println(jsonMsg)
		}
	}

}

func startserver() {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(message)
	}
}

func main() {

	flag.StringVar(&server, "server", "localhost", "server address (local is default)")
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

	//Setup channel for Coinbase Message
	ch := make(chan coinbasepro.Message)

	//Start Stream
	go stream(productIDList, ch)

	//Send Channel
	sender(ch)

	//Start Server (testing)
	go startserver()

}
