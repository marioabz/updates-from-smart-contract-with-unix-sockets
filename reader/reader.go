package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
)

type Event struct {
	Address string;
	Reserve0 *big.Int;
	Reserve1 *big.Int;
}

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
	event := Event{}
	for {
		fmt.Println("reading answer")
		var buff [1024]byte
		n, err := conn.Read(buff[:])
		abc := json.Unmarshal(buff[:n], &event)
		fmt.Println(abc)
		fmt.Printf("%s\n", string(buff[:n]));
		if err != nil {
			log.Fatal(err)
		}
	}
}
