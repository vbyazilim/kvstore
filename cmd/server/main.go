package main

import (
	"log"

	"github.com/vbyazilim/kvstore/src/apiserver"
)

func main() {
	if err := apiserver.New(); err != nil {
		log.Fatal(err)
	}
}
