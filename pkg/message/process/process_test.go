package process_test

import (
	"testing"

	"github.com/xfiendx4life/ponytest/pkg/message/process"
)

func TestRead(t *testing.T) {
	test := make(chan struct{}, 1)
	go func() {
		process.Work()
		test <- struct{}{}
	}()
	<-test

}
