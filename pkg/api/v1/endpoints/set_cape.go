package endpoints

import (
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// SetCape godoc
//
//	@Summary	Set a player's cape
//	@Tags		Cape
//	@Accept		json
//	@Produce	json
//	@Param		id		query	string	true	"Name of the cape to be set"
//	@Success	200		"Success"
//	@Failure	400		{object}	response.ValidationError	"Missing or invalid cape in query"
//	@Failure	404		{object}	response.Error				"Cape does not exist"
//	@Failure	500		{object}	response.Error				"Internal server error"
//	@Router		/cape 	[put]
//	@Security 	Bearer
func SetCape(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)
	player := ctx.MustGet("player").(response.Player)

	cape := ctx.Query("id")
	if cape == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Missing cape in query",
		})
		return
	}

	// Check if the cape has an entry
	_, err := cache.Get(ctx.Request.Context(), cape)
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting cape from cache", zap.String("cape", cape), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
		logger.Error("Client tried to set a cape that does not exist", zap.String("cape", cape), zap.Any("player", player))

		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "This cape does not exist.",
		})
		return
	}

	err = cache.Set(ctx.Request.Context(), player.UUID.String(), []byte(cape), 0, 0)
	if err != nil {
		logger.Error("Error setting cape", zap.String("id", player.UUID.String()), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}
