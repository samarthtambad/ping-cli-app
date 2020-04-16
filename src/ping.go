package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
)

const localAddr = "0.0.0.0"

func main() {

	// check usage
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}

	// resolve hostname
	remoteAddr, err := net.ResolveIPAddr("ip4", os.Args[1])
	if err != nil {
		fmt.Println("Resolution error", err.Error())
		os.Exit(1)
	}

	ping(remoteAddr)

}

func ping(remoteAddr *net.IPAddr)  {

	// establish connection
	conn, err := icmp.ListenPacket("ip4:icmp", localAddr)
	checkError(err)
	defer conn.Close()

	// prepare message
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code:     0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}

	// marshall packet
	msgBytes, err := msg.Marshal(nil)
	checkError(err)


	// send packet
	n, err := conn.WriteTo(msgBytes, remoteAddr)
	checkError(err)
	fmt.Print("Message sent: ", n, msgBytes)

	// receive a reply
	replyBytes := make([]byte, 1500)
	size, peer, err := conn.ReadFrom(replyBytes)
	checkError(err)

	recvMsg, err := icmp.ParseMessage(1, replyBytes[:size])
	checkError(err)

	fmt.Printf("Message received from %v: %d %v", peer, size, recvMsg.Type)
	
}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}