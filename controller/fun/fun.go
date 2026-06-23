package fun

import (
	"net/http"
	"strings"

	"luangao/biu"
	handlerfun "luangao/handler/fun"

	"github.com/gin-gonic/gin"
)

type RandomJumpController struct {
	randomJumpHandler handlerfun.RandomJumpFinder
}

func NewRandomJumpController(randomJumpHandler handlerfun.RandomJumpFinder) *RandomJumpController {
	if randomJumpHandler == nil {
		randomJumpHandler = handlerfun.NewRandomJumpHandler(nil)
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

func (c *RandomJumpController) GetThreeRandomJumps(ctx *gin.Context) {
	userID := extractUserID(ctx)
	results, err := c.randomJumpHandler.PickN(ctx.Request.Context(), 3, userID)
	if err != nil {
		biu.Failed(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	biu.Success(ctx, results)
}

func extractUserID(ctx *gin.Context) string {
	// Try JWT auth header first
	if id := biu.GetUserID(ctx); id != "" {
		return id
	}
	// Fall back to query param (for non-authenticated frontend calls)
	return strings.TrimSpace(ctx.Query("user"))
}
