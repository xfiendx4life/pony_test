package process

import (
	"context"
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xfiendx4life/ponytest/pkg/message/usecase"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type MQTTProcess struct {
	commonStorage *sync.Map
	worker        usecase.CaseWorker
	Done          chan *models.Message
	Ers           chan error
}

func New(common *sync.Map, pool usecase.CaseWorker) *MQTTProcess {
	return &MQTTProcess{
		commonStorage: common,
		worker:        pool,
		Done:          make(chan *models.Message, 1),
		Ers:           make(chan error, 1),
	}
}

func (mq *MQTTProcess) Work(ctx context.Context, host string, port int, stop chan struct{}) (err error) {
	var cnt int
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("MSG: %s\n", msg.Payload())
		go mq.worker(cnt, &msg, mq.Done, mq.Ers)
		cnt++
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")

	default:
		opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)) //"tcp://104.236.0.154:1883"
		opts.SetClientID("pony")
		opts.SetDefaultPublishHandler(f)
		topic := "/devices/+/state"

		opts.OnConnect = func(c mqtt.Client) {
			if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
				err = fmt.Errorf("stopped recieving messages: %s", token.Error())
			}
		}
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			err = fmt.Errorf("stopped recieving messages: %s", token.Error())
		} else {
			log.Printf("Connected to server\n")
		}
		<-stop
		err = fmt.Errorf("stopeed by user")
		log.Println(err)
		return err
	}
}
