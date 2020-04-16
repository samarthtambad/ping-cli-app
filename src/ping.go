package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
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

	duration, err := ping(remoteAddr, time.Now().Add(20 * time.Millisecond))
	if err != nil {  }
	fmt.Println(duration)

}

// send packet to remote address and receive response,
// return (success, duration)
func ping(remoteAddr *net.IPAddr, ttl time.Time) (time.Duration, error) {

	start := time.Now()

	// establish connection
	conn, err := icmp.ListenPacket("ip4:icmp", localAddr)
	checkError(err)
	defer conn.Close()

	// set ttl
	err = conn.SetDeadline(ttl)
	checkError(err)

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
	if err != nil { return 0, err}

	// send packet
	_, err = conn.WriteTo(msgBytes, remoteAddr)
	if err != nil { return 0, err}
	//fmt.Print("Message sent: ", n, msgBytes)

	// receive a reply
	replyBytes := make([]byte, 1500)
	size, _, err := conn.ReadFrom(replyBytes)
	if err != nil { return 0, err}

	duration := time.Since(start)

	recvMsg, err := icmp.ParseMessage(1, replyBytes[:size])
	if err != nil { return 0, err}

	//fmt.Printf("Message received from %v: %d %v", peer, size, recvMsg.Type)

	switch recvMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		return duration, nil

	default:
		return 0, fmt.Errorf("expected %s, got %s", ipv4.ICMPTypeEchoReply.String(), recvMsg.Type)
	}

}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}