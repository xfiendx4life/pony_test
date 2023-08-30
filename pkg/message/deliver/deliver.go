package deliver

import "github.com/labstack/echo/v4"

type Deliver interface {
	ListID(ctx echo.Context) error
	GetDataById(ctx echo.Context) error
	SendRPC(ctx echo.Context) error
}
