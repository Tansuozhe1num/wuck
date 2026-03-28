package httpserver

import (
	"net/http"

	controllerfun "luangao/controller/fun"
	handlerfun "luangao/handler/fun"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	return SetupRouterWithFinder(handlerfun.NewRandomJumpHandler())
}

func SetupRouterWithFinder(randomJumpFinder handlerfun.RandomJumpFinder) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	randomJumpController := controllerfun.NewRandomJumpController(randomJumpFinder)
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(homePage))
	})

	r.GET("/api/biu", randomJumpController.GetRandomJump)

	return r
}
