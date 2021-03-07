package main

import (
	"fmt"
	"log"
	"net"
	"os"
)


func main() {

	var socketAddress = "/tmp/reserves.sock"

	if err := os.RemoveAll(socketAddress); err != nil {
        log.Fatal(err)
	}
		
	_socket, err := net.Listen("unix", socketAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(socketAddress)
	conn, err := _socket.Accept()
	if err != nil {
		log.Fatal(err)
	}
	for {
		var buff [1024]byte
		n, err := conn.Read(buff[:])
		fmt.Printf("%s\n", string(buff[:n]));
		if err != nil {
			log.Fatal(err)
		}
	}
}
