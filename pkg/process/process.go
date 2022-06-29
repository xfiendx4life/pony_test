package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

type Payload struct {
	Method string `json:"method"`
	// TODO: change params to map
	Params []string `json:"params"`
}

func New(common *sync.Map, pool usecase.CaseWorker) *MQTTProcess {
	return &MQTTProcess{
		commonStorage: common,
		worker:        pool,
		Done:          make(chan *models.Message, 1),
		Ers:           make(chan error, 1),
	}
}

func (mq *MQTTProcess) Pub(ctx context.Context, client *mqtt.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, ok := mq.commonStorage.LoadAndDelete("rpc")
			// log.Println(data)
			if ok {
				d := data.(*models.Message)
				p := Payload{}
				err := json.NewDecoder(strings.NewReader(d.Data)).Decode(&p)
				if err != nil {
					// We lose message if it's not correct
					mq.Ers <- fmt.Errorf("can't parse payload: %s", err)
					continue
				}
				payload := struct {
					Ponyrpc  string   `json:"ponyrpc"`
					Method   string   `json:"method"`
					Params   []string `json:"params"`
					Id       int      `json:"id"`
					Retry_id int      `json:"retry_id"`
				}{
					Ponyrpc:  "2.0",
					Method:   p.Method,
					Params:   p.Params,
					Id:       0,
					Retry_id: 0,
				}
				pp, err := json.Marshal(payload)
				if err != nil {
					// We lose message if it's not correct
					mq.Ers <- fmt.Errorf("can't parse payload: %s", err)
					continue
				}
				// payload := fmt.Sprintf(`{"ponyrpc":"2.0","method":%s,"params":%v,"id":0,"retry_id":0}`,
				// 	p.Method, p.Params)
				token := (*client).Publish(fmt.Sprintf("/devices/%s/rpc", d.ID), 0, false, pp)
				log.Println(string(pp))
				if ok := token.Wait(); !ok && token.Error() != nil {
					mq.Ers <- fmt.Errorf("can't send to rpc: %s", token.Error())
				}
			}
		}
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
		go mq.Pub(ctx, &client)
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
