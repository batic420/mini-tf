package main

import (
	"fmt"
	"log"
	"os"

	"github.com/batic420/mini-tf/internal/creator"
	"github.com/batic420/mini-tf/internal/decoder"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: mini-tf <path-to-yaml>")
	}

	env, err := decoder.Load(os.Args[1])
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if error := creator.CreateResource(*env); error != nil {
		fmt.Println(error.Error())
	}
}
