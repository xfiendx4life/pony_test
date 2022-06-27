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
		pr.Work(ctx, "104.236.0.154", 1883, st)
	}()
	go func(ctx context.Context) {
		select {
		case m := <-pr.Done:
			// TODO: Write to sync.Mutex
			log.Println(m)
		case err := <-pr.Ers:
			log.Println(err)
			cancel()
		case <-ctx.Done():
			return
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
