package deliver

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type del struct {
	commonStorage *sync.Map
}

type RpcData struct {
	ID     string            `json:"id,omitempty"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"` // ? change to interface
}

type Payload struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

func New(cStorage *sync.Map) Deliver {
	return &del{
		commonStorage: cStorage,
	}
}

func (d *del) ListID(ctx echo.Context) error {
	res := make([]map[string]time.Time, 0)
	var err error
	d.commonStorage.Range(func(key, value any) bool {
		if key != "rpc" {
			v, ok := value.(*models.Message)
			if !ok {
				err = echo.NewHTTPError(http.StatusInternalServerError, "can't parse data")
			}
			res = append(res, map[string]time.Time{v.ID: v.TimeStamp})
		}
		return true
	})
	if err != nil {
		return err
	}
	if len(res) > 0 {
		return ctx.JSON(200, res)
	}
	return ctx.NoContent(http.StatusNoContent)

}

func (d *del) GetDataById(ctx echo.Context) error {
	id := ctx.Param("id")
	if data, ok := d.commonStorage.Load(id); ok {
		return ctx.JSON(http.StatusOK, data.(*models.Message))
	}
	return ctx.NoContent(http.StatusNotFound)

}
func (d *del) SendRPC(ctx echo.Context) error {
	data := RpcData{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		log.Printf("can't decode incoming json: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't decode data")
	}
	id := data.ID
	data.ID = ""
	pl, err := json.Marshal(data)
	if err != nil {
		log.Printf("can't encode incoming json: %s\n", err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't encode data")
	}
	fakeMessage := models.Message{
		ID:        id,
		Data:      string(pl),
		TimeStamp: time.Now(),
	}
	log.Printf("%#v", data.Params)
	d.commonStorage.Store("rpc", &fakeMessage)
	log.Println("data passed to storage")
	return ctx.NoContent(http.StatusOK)
}
