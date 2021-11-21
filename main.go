package main

import (
	"fmt"
	"log"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging/consolelogging/logformat"
	"github.com/eurofurence/reg-auth-service/web"
)

func main() {
	err := config.LoadConfiguration("config.yaml")
	if err != nil {
		log.Fatal(logformat.Logformat("ERROR", "00000000", fmt.Sprintf("Error while loading configuration: %v", err)))
	}

	database.Open()
	defer database.Close()

	server := web.Create()
	web.Serve(server)
}
