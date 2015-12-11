package handler

import "../interfaces"

// LoadHandlers returns all available handlers.
func LoadHandlers() []interfaces.Handler {
	return []interfaces.Handler{PingPongHandler{}, UranaiHandler{}}
}
