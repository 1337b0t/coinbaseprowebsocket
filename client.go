package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"github.com/preichenberger/go-coinbasepro"
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
