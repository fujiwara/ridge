package main

import (
	"context"
	"log"
	"os"

	"github.com/fujiwara/ridge/ridgecli"
)

func main() {
	if err := ridgecli.Run(context.TODO()); err != nil {
		log.Println("[error]", err)
		os.Exit(1)
	}
}
