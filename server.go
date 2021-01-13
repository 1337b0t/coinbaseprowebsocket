package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

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
