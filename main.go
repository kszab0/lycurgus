package main

import (
	"log"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatal(app.RunBlocker())
	}()
	app.RunGUI()
}
