package process_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/ponytest/pkg/message/process"
	"github.com/xfiendx4life/ponytest/pkg/message/usecase"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type mockStore struct {
	err error
}

func (mc *mockStore) Write(ctx context.Context, data models.Message) error {
	return mc.err
}

//TODO: check this test
func TestRead(t *testing.T) {
	uc := usecase.New(&mockStore{})
	pr := process.New(&sync.Map{}, uc.Produce())
	st := make(chan struct{})
	go func() {
		pr.Work(context.Background(), "104.236.0.154", 1883, st)
	}()
	for msg := range pr.Done {
		assert.NotNil(t, msg)
	}
	time.Sleep(20 * time.Second)
	st <- struct{}{}

}
