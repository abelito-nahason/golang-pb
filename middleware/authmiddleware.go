package middleware

import (
	"net/http"
	"pocketbase-seal/model/genericresponse"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func InitAuthMiddleware(se *router.RouterGroup[*core.RequestEvent]) {
	se.BindFunc(func(e *core.RequestEvent) error {
		if e.Request.Header.Get("Authorization") == "" {
			errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Missing auth token", ResponseCode: 400}
			return e.JSON(http.StatusBadRequest, errResponse)
		}
		return e.Next()
	})
}
