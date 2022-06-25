package process_test

import (
	"context"
	"testing"
	"time"

	"github.com/xfiendx4life/ponytest/pkg/message/process"
)

func TestRead(t *testing.T) {
	st := make(chan struct{})
	go func() {
		process.Work(context.Background(), "104.236.0.154", 1883, st)
	}()
	time.Sleep(20 * time.Second)
	st <- struct{}{}

}
