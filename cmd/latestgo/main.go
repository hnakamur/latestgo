package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hnakamur/latestgo"
)

func main() {
	ver, err := latestgo.Version(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(ver)
}
