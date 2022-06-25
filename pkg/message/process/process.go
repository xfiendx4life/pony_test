package process

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var knt int
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	// TODO: Concurrent procceed messages
	fmt.Printf("Topic: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	knt++
}

func Work() {
	knt = 0
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := MQTT.NewClientOptions().AddBroker("tcp://104.236.0.154:1883")
	opts.SetClientID("pony")
	opts.SetDefaultPublishHandler(f)
	topic := "/devices/+/state"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	<-c
}
