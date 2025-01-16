package main

import (
	"fmt"
	"message-app/rest"
	"message-app/services"
)

func main() {
	/*
	* Initiate Service Layer Container
	 */
	serviceContainer := services.NewServiceContainer()

	/*
	* Initiate Rest Server
	 */
	rest.StartServer(serviceContainer)

	fmt.Println("========== Rest Server Started ============")

	select {}
}
