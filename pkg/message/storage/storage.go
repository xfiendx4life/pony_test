package storage

import (
	"context"

	"github.com/xfiendx4life/ponytest/pkg/models"
)

type Storage interface {
	Write(ctx context.Context, message models.Message) error
}
