package main

import (
	"fmt"

	"github.com/samertm/sheep-mmo/server"
)

func main() {
	fmt.Println("Listening on ws://localhost:4977")
	server.ListenAndServe("localhost:4977")
}
