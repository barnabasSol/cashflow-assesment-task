package paymentmaking

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	group *echo.Echo
	s     Service
}

func InitHandler(s Service, group *echo.Group) *Handler {
	h := &Handler{
		s: s,
	}
	h.group.POST("/payment", h.CreatePayment)
	return h
}

func (h *Handler) CreatePayment(ctx echo.Context) error {
	var req CreatePaymentRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	return nil

}
