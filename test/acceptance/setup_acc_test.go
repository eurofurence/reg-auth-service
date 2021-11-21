package acceptance

import (
	"log"
	"net/http/httptest"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/web"
)

// placing these here because they are package global

var (
	ts *httptest.Server
)

const tstDefaultConfigFile = "../../test/resources/config-acceptancetests.yaml"

func tstSetup(configFilePath string) {
	tstSetupConfig(configFilePath)
	tstSetupHttpTestServer()
	database.Open()
}

func tstSetupConfig(configFilePath string) {
	err := config.LoadConfiguration(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func tstSetupHttpTestServer() {
	router := web.Create()
	ts = httptest.NewServer(router)
}

func tstShutdown() {
	database.Close()
	ts.Close()
}
