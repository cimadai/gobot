package handler

import "../interfaces"

func LoadHandlers() []interfaces.Handler {
	return []interfaces.Handler{PingPongHandler{}}
}
