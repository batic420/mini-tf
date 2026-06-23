package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	b, _ := json.MarshalIndent(env, "", " ")
	fmt.Println(string(b))
}
