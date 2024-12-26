package main

import (
	"log"
	"os"
	"pocketbase-seal/handler"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		sealGroup := se.Router.Group("/api/seal")
		sealHandler := handler.NewSealHandler(app)
		sealGroup.GET("/get/{id}", sealHandler.GetSpecificSeal)
		sealGroup.GET("/list", sealHandler.GetSeals)
		sealGroup.PUT("/update", sealHandler.UpdateSeal)
		sealGroup.POST("/add", sealHandler.CreateNewSeal)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
