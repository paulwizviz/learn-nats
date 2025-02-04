package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sub, _ := nc.Subscribe("greet.*", func(msg *nats.Msg) {
		name := msg.Subject[6:] // Extract subject from message
		msg.Respond([]byte("hello, " + name))
	})

	rep, _ := nc.Request("greet.joe", nil, time.Second)
	fmt.Println(string(rep.Data))

	rep, _ = nc.Request("greet.sue", nil, time.Second)
	fmt.Println(string(rep.Data))

	rep, _ = nc.Request("greet.bob", nil, time.Second)
	fmt.Println(string(rep.Data))

	sub.Unsubscribe()

	_, err = nc.Request("greet.joe", nil, time.Second)
	fmt.Println(err)

}
