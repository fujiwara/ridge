package main

import (
	"context"
	"log"
	"os"

	"github.com/fujiwara/ridge"
)

func main() {
	if err := ridge.RunCLI(context.TODO()); err != nil {
		log.Println("[error]", err)
		os.Exit(1)
	}
}
