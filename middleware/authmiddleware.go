package middleware

import (
	"log"
	"net/http"
	"pocketbase-seal/helper"
	"pocketbase-seal/model/auth"
	"pocketbase-seal/model/genericresponse"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func InitAuthMiddleware(se *router.RouterGroup[*core.RequestEvent], app *pocketbase.PocketBase) {
	se.BindFunc(func(e *core.RequestEvent) error {
		if e.Request.Header.Get("Authorization") == "" {
			errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Missing auth token", ResponseCode: 400}
			return e.JSON(http.StatusBadRequest, errResponse)
		}

		claims, err := helper.VerifyAuth(e.Request.Header.Get("Authorization"))
		if err != nil {
			log.Print(err)
			errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Error parsing token", ResponseCode: 400}
			return e.JSON(http.StatusBadRequest, errResponse)
		}

		splitUrl := strings.Split(e.Request.Pattern, " ")

		rbac := auth.RoleRoute{}

		err = app.DB().NewQuery("SELECT accessible FROM roleroutes INNER JOIN routes ON routes.id = roleroutes.route INNER JOIN roles ON roles.id = roleroutes.role WHERE roles.name = {:role} AND routes.route = {:route}").Bind(dbx.Params{
			"role":  claims.Role,
			"route": splitUrl[1],
		}).One(&rbac)

		if err != nil {
			errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Error receiving rbac", ResponseCode: 500}
			if strings.Contains(err.Error(), "no rows in result") {
				errResponse = genericresponse.GenericErrorResponse{ResponseMessage: "Routes/roles combination is not yet initialized", ResponseCode: 404}
				return e.JSON(http.StatusNotFound, errResponse)
			}
			return e.JSON(http.StatusInternalServerError, errResponse)
		}

		if !rbac.Accessible {
			errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Unauthorized", ResponseCode: 401}
			return e.JSON(http.StatusUnauthorized, errResponse)
		}

		return e.Next()
	})
}
