package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/xfiendx4life/ponytest/pkg/message/deliver"
	"github.com/xfiendx4life/ponytest/pkg/message/storage"
	"github.com/xfiendx4life/ponytest/pkg/message/usecase"
	"github.com/xfiendx4life/ponytest/pkg/process"
	"github.com/xfiendx4life/ponytest/pkg/rest"
)

func main() {
	var commonStorage = sync.Map{}
	store, err := storage.New("storage")
	if err != nil {
		log.Fatalf("can't create store %s", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	uc := usecase.New(store)
	pr := process.New(&commonStorage, uc.Produce())
	del := deliver.New(&commonStorage)
	st := make(chan struct{})
	server := rest.New(del)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	go func() {
		if er := server.StartServer(ctx, host, port); er != nil {
			log.Fatalf("server stopped %s", err)
		}
	}()
	go func() {
		err = pr.Work(ctx, "104.236.0.154", 1883, st)
		log.Println(err)
	}()
	go func(ctx context.Context) {
		for {
			select {
			case m := <-pr.Done:
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
