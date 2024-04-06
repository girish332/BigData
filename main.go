package main

import (
	routerPkg "github.com/girish332/bigdata/router"
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
