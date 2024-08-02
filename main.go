package main

import (
	"log"

	"event-service/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Panic(r)
		}
	}()

	cmd.Execute()
}
