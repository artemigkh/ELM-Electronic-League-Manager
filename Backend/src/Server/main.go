package main

import (
	"Server/config"
	"Server/routes"
	_ "github.com/jackc/pgx/stdlib"
	"log"
)

func main() {
	conf := config.GetConfig("Backend/conf.json")

	//start router/webapp
	app := routes.InitRoutes(conf)
	err := app.Run(conf.GetPortString())
	if err != nil {
		log.Fatal(err)
	}
}
