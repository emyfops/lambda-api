package endpoints

import (
	"github.com/Edouard127/lambda-api/internal/app/auth"
	"github.com/Edouard127/lambda-api/internal/app/memory"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	partyCountTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lambda_rpc_party_count_total",
		Help: "Total number of parties",
	}, []string{"version"})
)

func init() {
	prometheus.MustRegister(partyCountTotal)
}

// Player -> &Party ID
var playerMap = memory.NewCache[response.Player, uuid.UUID]()

// Reverse mapping of playerMap
// &Party ID -> &Party
var partyMap = memory.NewCache[uuid.UUID, *response.Party]()

// CreateParty 	godoc
//
//	@Summary	Create a new party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Param		settings	body		request.Settings	false	"Party configuration"
//	@Success	201			{object}	response.Party
//	@Failure	400			{object}	response.ValidationError
//	@Failure	409			{object}	response.Error
//	@Router		/party/create [post]
//	@Security	ApiKeyAuth
func CreateParty(ctx *gin.Context) {
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	_, exists := playerMap.Get(player)
	if exists {
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Error{
			Message: "You are already in a party",
		})
		return
	}

	party := response.NewParty(player, &settings)

	partyMap.Set(party.ID, party, memory.NoExpiration)
	playerMap.Set(player, party.ID, memory.NoExpiration)

	partyCountTotal.WithLabelValues("v1").Inc()

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
}
