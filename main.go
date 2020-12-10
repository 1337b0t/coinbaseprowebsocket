package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
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

	for true {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			panic("Script restart required")
		}

		ch <- message

	}

}

func tcpclient(ch <-chan coinbasepro.Message) {

	//setup conn for tcp client
	conn, _ = net.Dial("tcp", serverAddr)

	for {

		select {
		case v, err := <-ch:
			if !err {
				return
			}

			//Skip blank price messages
			if v.Price != "" {
				tm := strconv.Itoa(int(v.Time.Time().Unix()))
				jsonMsg := fmt.Sprintf("[{market:%s,price:%s,size:%s,side: %s, bestBid:%s, bestAsk:%s, tm:%s}]", v.ProductID, v.Price, v.LastSize, v.Side, v.BestBid, v.BestAsk, tm)
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
		//fmt.Print("Message Received:", string(message))
		fmt.Print(string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		connServer.Write([]byte(newmessage + "\n"))
	}
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
	tcpclient(ch)

}
