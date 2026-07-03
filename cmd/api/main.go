package main

import (
	"log"

	"payment_gateway/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
