package main

import (
	"log"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.RunBlocker())
}
