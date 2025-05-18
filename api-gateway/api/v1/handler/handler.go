package handler

import (
	"api-gateway/service"
)

type Handler struct {
	service    service.IClients
}

type HandlerCinfig struct {
	Service service.IClients
}

func NewHandler(h *HandlerCinfig) *Handler {
	return &Handler{
		service: h.Service,
	}
}
