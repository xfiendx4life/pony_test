package deliver

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/ponytest/pkg/models"
)

type del struct {
	commonStorage *sync.Map
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
	return nil
}
