package main

import (
	"log"
	"time"
)

func main() {
	for {
		log.Println("lol")
		time.Sleep(time.Duration(1)*time.Minute)
	}
}

