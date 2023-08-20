package main

import (
	"log"
	"os"

	"github.com/vbyazilim/kvstore/src/apiserver"
)

func main() {
	if err := apiserver.New(
		apiserver.WithServerEnv(os.Getenv("SERVER_ENV")),
	); err != nil {
		log.Fatal(err)
	}
}
