package fun

import (
	"net/http"

	"luangao/biu"
	handlerfun "luangao/handler/fun"

	"github.com/gin-gonic/gin"
)

type RandomJumpController struct {
	randomJumpHandler handlerfun.RandomJumpFinder
}

func NewRandomJumpController(randomJumpHandler handlerfun.RandomJumpFinder) *RandomJumpController {
	if randomJumpHandler == nil {
		randomJumpHandler = handlerfun.NewRandomJumpHandler()
	}

	return &RandomJumpController{
		randomJumpHandler: randomJumpHandler,
	}
}

func (c *RandomJumpController) GetRandomJump(ctx *gin.Context) {
	result, err := c.randomJumpHandler.Pick(ctx.Request.Context())
	if err != nil {
		biu.Failed(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	biu.Success(ctx, result)
}
