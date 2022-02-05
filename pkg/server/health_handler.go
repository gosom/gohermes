package server

import (
	"net/http"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/utils"
)

func HealthHandler(di *container.ServiceContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.RenderJson(r, w, http.StatusOK, map[string]bool{
			"ready": true,
		})

	}
}
