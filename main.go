package main

import (
	"log"
	"os"
	"pocketbase-seal/handler"
	"pocketbase-seal/middleware"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env")
	}

	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		authHandler := handler.NewAuthHandler(app)
		authGroup := se.Router.Group("/auth")
		authGroup.POST("/register", authHandler.RegisterNewUser)
		authGroup.POST("/login", authHandler.LoginUser)

		apiGroup := se.Router.Group("/api")
		middleware.InitAuthMiddleware(apiGroup)

		sealGroup := apiGroup.Group("/seal")
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
