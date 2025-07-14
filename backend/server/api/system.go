package api

import (
	"meme3/service/model"
	"meme3/service/store"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"go.mongodb.org/mongo-driver/bson"
)

func initSystem(api *swag.API) {
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/system/banner",
			endpoint.Tags("System"),
			endpoint.Handler(webSystemBanner),
			endpoint.Summary("system banner"),
			endpoint.Description("get system banner"),
			endpoint.Response(http.StatusOK, "successed", endpoint.SchemaResponseOption([]string{})),
			// endpoint.Security("petstore_auth", "read:pets", "write:pets"),
		),
	)
}

func webSystemBanner(w http.ResponseWriter, r *http.Request) {
	if system, err := store.DBGet[model.System](model.C_System, bson.M{}); err != nil {
		WebResponseJson(w, r, ApiError(err.Error()), http.StatusOK)
	} else {
		WebResponseJson(w, r, ApiResponse(system.Banner, true), http.StatusOK)
	}
}
