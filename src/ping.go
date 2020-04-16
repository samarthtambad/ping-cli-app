package main

import (
	"fmt"
	"os"
)

func main() {

	// check usage
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}

}