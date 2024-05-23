package handler

import "38hw/storage"

type handler struct {
	storage storage.IStorage
}

type HandlerConfig struct {
	Storage storage.IStorage
}

func New(c *HandlerConfig) *handler {
	return &handler{
		storage: c.Storage,
	}
}
