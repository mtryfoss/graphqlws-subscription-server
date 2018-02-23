package gss

import gss "github.com/taiyoh/graphqlws-subscription-server"

type Handler struct {
	listener *gss.Listener
}

func NewHandler(listener *gss.Listener) *Handler {
	return &Handler{
		listener: listener,
	}
}
