package handler

import (
	"api-gateway/service"
)

type Handler struct {
	service    service.IClients
}

type HandlerConfig struct {
	Service service.IClients
}

func NewHandler(h *HandlerConfig) *Handler {
	return &Handler{
		service: h.Service,
	}
}