package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	// read input
	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading stdin: %v", err)
	}
	json.NewEncoder(os.Stdout).Encode(struct {
		Body string
	}{
		Body: string(body),
	})
}
