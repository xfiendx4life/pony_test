package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xfiendx4life/ponytest/pkg/message/storage"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type uCase struct {
	store storage.Storage
}

func New(st storage.Storage) Usecase {
	return &uCase{
		store: st,
	}
}

func (u *uCase) Produce() CaseWorker {
	res := func(id int, data *mqtt.Message, done chan *models.Message, err chan error) {
		log.Printf("worker %d started", id)
		log.Printf("Topic: %s\n", (*data).Topic())
		name := strings.Split((*data).Topic(), "/")[2]
		res := models.Message{
			ID:        name,
			Data:      string((*data).Payload()),
			TimeStamp: time.Now(),
		}
		go func() {
			if er := u.store.Write(context.Background(), res); er != nil {
				fmt.Println("from wrkr")
				err <- er
			}
		}()
		done <- &res
		log.Printf("wrkr %d sent result\n", id)
	}

	return res
}
