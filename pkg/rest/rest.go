package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/ponytest/pkg/message/deliver"
)

type RestServer struct {
	del deliver.Deliver
	*echo.Echo
}

func New(del deliver.Deliver) *RestServer {
	return &RestServer{
		del:  del,
		Echo: echo.New(),
	}
}

func (r *RestServer) StartServer(ctx context.Context, host string, port int) (err error) {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		r.GET("/:id", r.del.GetDataById)
		r.GET("/list", r.del.ListID)
		go func() {
			log.Printf("starting server at %s:%d\n", host, port)
			r.HideBanner = true
			if err = r.Start(fmt.Sprintf("%s:%d", host, port)); err != nil {
				log.Printf("error while working with server: %s\n", err)
			}
		}()
	}
	return err
}
