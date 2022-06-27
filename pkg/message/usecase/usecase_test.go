package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/ponytest/pkg/message/usecase"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type mockStore struct {
	err error
}

func (mc *mockStore) Write(ctx context.Context, data models.Message) error {
	return mc.err
}

// TODO: More tests for usecase
func TestProduce(t *testing.T) {
	u := usecase.New(&mockStore{})
	cw := u.Produce()
	require.NotNil(t, cw)
}
