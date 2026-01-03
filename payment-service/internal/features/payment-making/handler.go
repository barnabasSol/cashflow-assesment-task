package paymentmaking

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	e *echo.Echo
	s Service
}

func InitHandler(s Service, e *echo.Echo) *Handler {
	h := &Handler{
		s: s,
		e: e,
	}
	h.e.POST("/payment", h.CreatePayment)
	return h
}

func (h *Handler) CreatePayment(ctx echo.Context) error {
	var req CreatePaymentRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	res, err := h.s.CreatePayment(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusCreated, res)
	return nil
}
