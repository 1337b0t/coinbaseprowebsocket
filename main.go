package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

var server, port, productID string
var serverAddr string
var productIDList []string
var priceSizeSum, sizeSum float64
var conn net.Conn

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

func tcpclient(ch <-chan coinbasepro.Message) {

	for {

		select {
		case v, err := <-ch:
			if !err {
				return
			}

			//Skip blank price messages
			if v.Price != "" {
				//Message:
				//Market: BTC-USD
				//Price:  19000
				//Size:   1
				priceFloat, _ := strconv.ParseFloat(v.Price, 64)
				sizeFloat, _ := strconv.ParseFloat(v.LastSize, 64)

				priceSizeSum += priceFloat * sizeFloat
				sizeSum += sizeFloat
				vwap := priceSizeSum / sizeSum

				jsonMsg := fmt.Sprintf("[{market:%s,price:%s,vwap:%.2f,size:%s}]", v.ProductID, v.Price, vwap, v.LastSize)
				// send to socket
				fmt.Fprintf(conn, jsonMsg+"\n")
				// listen for reply
				bufio.NewReader(conn).ReadString('\n')
				//fmt.Print("Message from server: " + message)
			}

		}

	}
}

func tcpserver() {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	connServer, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(connServer).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		connServer.Write([]byte(newmessage + "\n"))
	}
}

func main() {

	flag.StringVar(&server, "server", "localhost", "server address (local is default)")
	flag.StringVar(&port, "port", "8081", "the default port number is 8081")
	flag.StringVar(&productID, "products", "BTC-USD", "BTC-USD")
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
	//setup conn for tcp client
	conn, _ = net.Dial("tcp", serverAddr)
	tcpclient(ch)

}
