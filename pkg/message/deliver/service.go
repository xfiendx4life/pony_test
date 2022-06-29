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
	ID     string   `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"` // ? change to interface
}

func New(cStorage *sync.Map) Deliver {
	return &del{
		commonStorage: cStorage,
	}
}

func (d *del) ListID(ctx echo.Context) error {
	res := make([]map[string]time.Time, 0)
	d.commonStorage.Range(func(key, value any) bool {
		v := value.(*models.Message)
		res = append(res, map[string]time.Time{v.ID: v.TimeStamp})
		return true
	})
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
	d.commonStorage.Store("rpc", data)
	log.Println("data passed to storage")
	return ctx.NoContent(http.StatusOK)
}
