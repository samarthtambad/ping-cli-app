package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const localAddr = "0.0.0.0"

func main() {

	// check usage
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Usage: program", "[-ttl]", "[-ipv6]", "host")
		os.Exit(1)
	}

	// parse optional ttl flag
	ttlPtr := flag.Int("ttl", 255, "Set the IP Time To Live for outgoing packets")
	ipv6Ptr := flag.Bool("ipv6", false, "Set the protocol to IPv6")
	flag.Parse()
	fmt.Printf("ttl: %d, ipv6: %v, host: %s \n", *ttlPtr, *ipv6Ptr, flag.Arg(0))

	// set protocol dependent values
	var network string
	if *ipv6Ptr {		// if ipv6 flag is set
		network = "ip6"
	} else {
		network = "ip4"
	}

	// resolve hostname
	remoteAddr, err := net.ResolveIPAddr(network, flag.Arg(0))
	if err != nil {
		fmt.Println("Resolution error", err.Error())
		os.Exit(1)
	}

	// notify on exit interrupt (^C)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	pingInterval := 2 * time.Second
	pingTicker := time.NewTicker(pingInterval)

	// periodically send echo requests
	seqNo := 0
	recv := 0
	go func(ttl int) {
		fmt.Printf("PING %s \n", remoteAddr.String())
		for {
			select {
			case <-pingTicker.C:
				duration, err := ping(remoteAddr, ttl, *ipv6Ptr)
				if err != nil {
					fmt.Printf("Request timeout for icmp_seq %d no route to host %s \n", seqNo, remoteAddr.String())
				} else {
					fmt.Printf("Response from %s: icmp_seq=%d packets_lost=%d ttl=%d latency=%v \n", remoteAddr.String(), seqNo, seqNo - recv, ttl, duration.String())
					recv += 1
				}
				seqNo += 1
			}
		}
	}(*ttlPtr)

	<-done
	pingTicker.Stop()
	fmt.Printf("--- %s ping statistics ---\n", remoteAddr.String())
	fmt.Printf("%d packets transmitted, %d packets received, %d%% packet loss \n", seqNo, recv, ((seqNo-recv) * 100)/seqNo)
}

// send packet to remote address and receive response,
// return (success, duration)
func ping(remoteAddr *net.IPAddr, ttl int, v6 bool) (time.Duration, error) {

	start := time.Now()

	// set protocol dependent values
	var network string
	var msgType icmp.Type
	var proto int
	if v6 {		// if ipv6 flag is set
		network = "ip6:icmp"
		msgType = ipv6.ICMPTypeEchoRequest
		proto = 1
	} else {
		network = "ip4:icmp"
		msgType = ipv4.ICMPTypeEcho
		proto = 1
	}

	// establish connection
	conn, err := icmp.ListenPacket(network, localAddr)
	if err != nil { return 0, err}
	defer conn.Close()

	// set TTL/hop limit
	if v6 {
		conn.IPv6PacketConn().SetHopLimit(ttl)
	} else {
		conn.IPv4PacketConn().SetTTL(ttl)
	}

	// set deadline of 1s to limit indefinite wait for response
	conn.SetDeadline(time.Now().Add(1 * time.Second))

	// prepare message
	msg := icmp.Message{
		Type: msgType,
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

	// receive a reply
	replyBytes := make([]byte, 1500)
	size, _, err := conn.ReadFrom(replyBytes)
	if err != nil { return 0, err}

	duration := time.Since(start)

	recvMsg, err := icmp.ParseMessage(proto, replyBytes[:size])
	if err != nil { return 0, err}

	// check received message type
	switch recvMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		return duration, nil
	case ipv6.ICMPTypeEchoReply:
		return duration, nil
	default:
		return 0, fmt.Errorf("expected %s or %s, got %s", ipv4.ICMPTypeEchoReply.String(), ipv6.ICMPTypeEchoReply.String(), recvMsg.Type)
	}

}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}