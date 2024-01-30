package main

import (
	routerPkg "BigData/router"
	"log"
)

func main() {
	router := routerPkg.InitializeRouter()
	err := router.Run(":8080")
	if err != nil {
		log.Println("error starting server: %v", err)
		return
	}
}
