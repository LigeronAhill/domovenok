package main

import (
	"log"
	"os"
	"strings"

	"github.com/LigeronAhill/domovenok/internal/config"
)

func main() {
	env := "local"
	if prod := os.Getenv("ENVIRONMENT"); prod != "" {
		if strings.ToLower(prod) == "production" {
			env = "production"
		}
	}
	_, err := config.Init(env)
	if err != nil {
		log.Fatal(err)
	}
}
