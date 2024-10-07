package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"linkshortener/innards"
)

func main() {
	var r = gin.Default()

	r.HEAD("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	innards.AuthEndpoints(r)
	innards.ApiEndpoints(r)
	innards.UrlMapping(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
