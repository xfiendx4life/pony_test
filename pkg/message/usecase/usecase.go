package usecase

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type CaseWorker = func(id int, job *mqtt.Message, done chan *models.Message, err chan error)

type Usecase interface {
	Produce() CaseWorker
}
