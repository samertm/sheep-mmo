package main

import (
	"fmt"
	"flag"
	"strconv"

	"github.com/samertm/sheep-mmo/server"
)

func main() {
	hostname := flag.String("hostname", "localhost", "the hostname")
	port := flag.Int("port", 4977, "the port")
	flag.Parse()
	fmt.Printf("Listening on ws://%s:%d", *hostname, *port)
	server.ListenAndServe((*hostname) + ":" + strconv.Itoa(*port))
}
