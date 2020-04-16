package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"os"
)

func main() {

	// check usage
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}

	const localAddr = "0.0.0.0"

	// establish connection
	conn, err := icmp.ListenPacket("ip4:icmp", localAddr)
	checkError(err)
	defer conn.Close()

}

func checkError(err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}