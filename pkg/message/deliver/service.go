package deliver

import (
	"net/http"
	"sync"

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
	res := make([]models.Message, 0)
	d.commonStorage.Range(func(key, value any) bool {
		res = append(res, value.(models.Message))
		return true
	})
	if len(res) > 0 {
		return ctx.JSON(200, res)
	}
	return ctx.NoContent(http.StatusNoContent)

}

// TODO: Test this shit
func (d *del) GetDataById(ctx echo.Context) error {
	id := ctx.Param("id")
	if data, ok := d.commonStorage.Load(id); ok {
		return ctx.JSON(http.StatusOK, data.(models.Message))
	}
	return ctx.NoContent(http.StatusNotFound)

}
func (d *del) SendRPC(ctx echo.Context) error {
	return nil
}
