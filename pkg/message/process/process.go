package process

import (
	"context"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Creates pool of workers to proceed messages
// func createPool(num int, )
// //TODO: set num in config

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	// TODO: Concurrent procceed messages

	fmt.Printf("Topic: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func Work(ctx context.Context, host string, port int, stop chan struct{}) (err error) {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)) //"tcp://104.236.0.154:1883"
		opts.SetClientID("pony")
		opts.SetDefaultPublishHandler(f)
		topic := "/devices/+/state"

		opts.OnConnect = func(c MQTT.Client) {
			if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
				err = fmt.Errorf("stopped recieving messages: %s", token.Error())
			}
		}
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			err = fmt.Errorf("stopped recieving messages: %s", token.Error())
		} else {
			log.Printf("Connected to server\n")
		}
		<-stop
		return fmt.Errorf("stopeed by user")
	}

}
