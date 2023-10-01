package main

import (
	"github.com/aswinjithkukku/unknownus-chatapp/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	hub := ws.NewHub()
	handler := ws.NewHandler(hub)
	go hub.Run()
	r.GET("/ws", handler.HandleWS)

	r.Run(":8080")
}
