package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/xfiendx4life/ponytest/pkg/models"
)

type fileSystem struct {
	path string
	mu   sync.Mutex
}

// New storage in filesystem
// path to dir as parameter
func New(path string) (Storage, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("can't get fullpat: %s", err)
	}
	_, err = os.ReadDir(fullPath)
	if os.IsNotExist(err) {
		log.Printf("can't open dir %s\n", err)
		err = os.Mkdir(fullPath, 0777)
		if err != nil {
			return nil, fmt.Errorf("can't create dir: %s", err)
		}
		log.Printf("dir %s created \n", fullPath)
	}
	return &fileSystem{
		path: fullPath,
	}, nil
}

//TODO: change it to proceed wrong data
func parseMessage(data models.Message) []byte {
	return []byte(fmt.Sprintf("%s: %s - %s\n", data.TimeStamp.Format("2006/01/02-03:04:05"), data.ID, data.Data))
}

// Write is a method to write data to file
func (fs *fileSystem) Write(ctx context.Context, data models.Message) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		fs.mu.Lock()
		defer fs.mu.Unlock()
		filename := filepath.Join(fs.path, "storage")
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			if file, err = os.Create(filename); err != nil {
				return fmt.Errorf("can't open file: %s", err)
			}
		}
		defer func() {
			err = file.Close()
		}()
		if _, err = file.Write(parseMessage(data)); err != nil {
			return fmt.Errorf("can't write data: %s", err)
		}
		return err
	}
}
