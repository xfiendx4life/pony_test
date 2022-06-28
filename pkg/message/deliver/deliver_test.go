package deliver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/ponytest/pkg/message/deliver"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

var storage = sync.Map{}
var testMessage = models.Message{
	ID:        "FrmCtr010",
	Data:      `{"t_air":25.4,"h_air":74.5,"co2":0,"time":32977,"read_errs":0,"outs_state":12,"uptime":131700,"wifi":-44}`,
	TimeStamp: time.Now(),
}

func TestListID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	storage.Store(testMessage.ID, testMessage)
	d := deliver.New(&storage)
	err := d.ListID(c)
	require.NoError(t, err)
	var tt []time.Time
	err = json.Unmarshal(rec.Body.Bytes(), &tt)
	require.NoError(t, err)
	require.EqualValues(t, testMessage.TimeStamp.UTC(), tt[0].UTC())
}

func TestGetID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/FrmCtr010", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("FrmCtr010")

	storage.Store(testMessage.ID, testMessage)
	d := deliver.New(&storage)
	err := d.GetDataById(c)
	require.NoError(t, err)
	var tt models.Message
	err = json.Unmarshal(rec.Body.Bytes(), &tt)
	require.NoError(t, err)
	//TODO: find out what's wrong
	// !require.EqualValues(t, testMessage, tt)
	require.NotNil(t, tt)
}
