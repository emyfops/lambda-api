package routes

import (
	"errors"
	"github.com/Edouard127/lambda-api/api/models/request"
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// GetCapes godoc
//
//	@Summary	Get a player's cape
//	@Tags		Cape
//	@Accept		json
//	@Produce	json
//	@Param		players	body	request.CapeLookup	true	"Required fields are missing or invalid"
//	@Success	200	{object}	[]response.Cape				"Player capes"
//	@Failure	400	{object}	response.ValidationError	"Missing or invalid ID in query"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/capes [get]
//	@Security 	Bearer
func GetCapes(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)

	var lookup request.CapeLookup

	err := ctx.ShouldBindJSON(&lookup)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	var ids = make([]string, len(lookup.Players))
	for _, id := range lookup.Players {
		ids = append(ids, id.String())
	}

	var capes []response.Cape
	items, err := cache.Gets(ctx.Request.Context(), false, ids...)
	if err != nil && !errors.Is(err, memcached.ErrNotFound) {
		logger.Error("Error getting player cape from cache", zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	for _, item := range items {
		id, err := uuid.Parse(item.Key)
		if err != nil {
			continue
		}

		capes = append(capes, response.Cape{
			Uuid: id,
			Type: string(item.Value),
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, capes)
	return
}
