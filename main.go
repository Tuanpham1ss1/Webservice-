package main

import (
	"log"
	"test1/infrastructure"
	"test1/router"
)

func main() {
	app := router.Router()
	log.Println("Database name: ", infrastructure.GetDBName())
	log.Printf("Server running at port: %+v\n", infrastructure.GetAppPort())
	log.Fatal(app.Listen(":" + infrastructure.GetAppPort()))
}
