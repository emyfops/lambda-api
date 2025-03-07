package endpoints

import (
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// GetCape godoc
//
//	@Summary	Get a player's cape
//	@Tags		Cape
//	@Accept		json
//	@Produce	json
//	@Param		id	query	string	true	"Player's ID"
//	@Success	200	{object}	response.Cape				"Player cape url"
//	@Failure	400	{object}	response.ValidationError	"Missing or invalid ID in query"
//	@Failure	404	{object}	response.Error				"No cape found for the provided ID"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/cape [get]
//	@Security 	Bearer
func GetCape(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)

	id := ctx.Query("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Missing ID in query",
		})
		return
	}

	item, err := cache.Get(ctx.Request.Context(), id)
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting player cape from cache", zap.String("id", id), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "This player does not have a cape.",
		})
		return
	}

	capeId := string(item.Value)

	item, err = cache.Get(ctx.Request.Context(), capeId)
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting cape from cache", zap.String("id", id), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
		logger.Error("Client tried to get a cape that does not exist", zap.String("id", id))

		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "This player does not have a cape.",
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Cape{
		Url:  string(item.Value),
		Type: capeId,
	})
	return
}
