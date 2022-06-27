package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/xfiendx4life/ponytest/pkg/message/process"
	"github.com/xfiendx4life/ponytest/pkg/message/storage"
	"github.com/xfiendx4life/ponytest/pkg/message/usecase"
)

var commonStorage = sync.Map{}

func main() {
	store, err := storage.New("storage")
	if err != nil {
		log.Fatalf("can't create store %s", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	uc := usecase.New(store)
	pr := process.New(&sync.Map{}, uc.Produce())
	st := make(chan struct{})
	go func() {
		err = pr.Work(ctx, "104.236.0.154", 1883, st)
		log.Println(err)
	}()
	go func(ctx context.Context) {
		for {
			select {
			case m := <-pr.Done:
				// TODO: Write to sync.Mutex
				commonStorage.Store(m.ID, m)
				log.Println(m)
			case err := <-pr.Ers:
				log.Println(err)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			cancel()
		}
	}
}
