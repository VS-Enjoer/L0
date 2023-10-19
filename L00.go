package main

import (
	"L00/Cache"
	"L00/DB"
	"L00/DB/Model"
	"L00/HTML"
	"L00/Nats"
	"github.com/nats-io/stan.go"
	"log"
)

var cache = make(map[string]Model.ModelOrder)

func main() {

	dbOpen, err := DB.DbConnect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return
	}
	go func() {
		if err := Cache.LoadCacheFromDB(dbOpen, cache); err != nil {
			log.Printf("Failed to load cache from DB: %v", err)
		}
	}()

	natsConn, err := Nats.СonnectToNats()

	defer natsConn.Close()

	_, err = natsConn.Subscribe("L00-channel", func(msg *stan.Msg) {
		go Nats.MessageHandler(msg, dbOpen, cache)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS: %v", err)
	}

	log.Print("Vse Good")

	go func() {
		HTML.HtmlStart(cache)
	}()

	select {}
}
