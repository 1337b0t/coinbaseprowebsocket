package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/preichenberger/go-coinbasepro"
)

const (
	pongWait = 60 * time.Second
)

func tcpclient(ch <-chan coinbasepro.Message, serverAddr string) {

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
				jsonMsg := fmt.Sprintf("[{market:%s,price:%s,size:%s,side:%s,bestBid:%s,bestAsk:%s,tm:%s}]", v.ProductID, v.Price, v.LastSize, v.Side, v.BestBid, v.BestAsk, tm)
				// send to socket
				fmt.Fprintf(conn, jsonMsg+"\n")
				// listen for reply
				bufio.NewReader(conn).ReadString('\n')
				//fmt.Print("Message from server: " + message)
			}

		}

	}
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Now().Sub(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}
