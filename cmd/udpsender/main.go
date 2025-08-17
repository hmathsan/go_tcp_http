package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	sender, err := net.DialUDP("udp", nil, conn)
	if err != nil {
		log.Fatal(err)
	}

	defer sender.Close()

	buffer := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		message, err := buffer.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		sender.Write([]byte(message))
	}
}
