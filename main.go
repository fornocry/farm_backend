package main

import (
	"crazyfarmbackend/config"
	"crazyfarmbackend/config/di"
	"crazyfarmbackend/src/api"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	port := "8000"

	init := di.Init()
	app := api.Init(init)
	app.Run(":" + port)
}
