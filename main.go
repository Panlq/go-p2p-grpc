package main

import (
	"log"

	"github/panlq-github/go-p2p-grpc/cmd"
)

func main() {
	// TODO: implement main
	if err := cmd.NewP2PCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
