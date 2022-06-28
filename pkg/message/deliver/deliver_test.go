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
	testResponse, _ := json.Marshal([]models.Message{testMessage})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/FrmCtr010", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	storage.Store(testMessage.ID, testMessage)
	d := deliver.New(&storage)
	err := d.ListID(c)
	require.NoError(t, err)

	require.EqualValues(t, string(testResponse), rec.Body.String()[:rec.Body.Len()-1])
}
