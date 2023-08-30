package storage_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/ponytest/pkg/message/storage"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

var testData = models.Message{
	ID:        "testId",
	TimeStamp: time.Now(),
	Data:      `{"time":51354,"uptime":922800,"rssi":-58,"cfghsh":847,"outs":[{"n":0,"s":0,"rt":0},{"n":1,"s":1,"rt":-1}],"override":[0,-1]}`,
}

var valid = `testId - {"time":51354,"uptime":922800,"rssi":-58,"cfghsh":847,"outs":[{"n":0,"s":0,"rt":0},{"n":1,"s":1,"rt":-1}],"override":[0,-1]}`

func checkFile() string {
	p, _ := filepath.Abs("./testpath")
	p = filepath.Join(p, "storage")
	data, err := os.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	strs := strings.Split(string(data), "\n")
	res := strings.Split(strs[len(strs)-2], ": ")
	return res[1]
}

func clean(path string) {
	p, _ := filepath.Abs(path)
	os.RemoveAll(p)
}

func TestWrite(t *testing.T) {
	st, err := storage.New("./testpath")
	require.NoError(t, err)
	defer clean("./testpath")
	err = st.Write(context.Background(), testData)
	require.NoError(t, err)
	require.Equal(t, valid, checkFile())
}

func TestWriteError(t *testing.T) {
	path := "/sys"
	st, err := storage.New(path)
	require.NoError(t, err)
	defer clean(path)
	err = st.Write(context.Background(), testData)
	require.Error(t, err)
}

func TestWriteContext(t *testing.T) {
	st, err := storage.New("./testpath")
	require.NoError(t, err)
	defer clean("./testpath")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d := make(chan struct{}, 1)
	go func() {
		err = st.Write(ctx, testData)
		d <- struct{}{}
	}()
	<-d
	require.Error(t, err)
	require.Equal(t, "done with context", err.Error())
}

// TODO: More tests here
// TODO: concurrent writing
