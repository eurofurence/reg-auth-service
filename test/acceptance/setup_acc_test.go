package acceptance

import (
	"context"
	"log"
	"net/http/httptest"
	"time"

	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/database"
	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/eurofurence/reg-auth-service/web"
	"github.com/eurofurence/reg-auth-service/web/controller/dropoffctl"
)

// placing these here because they are package global

var (
	ts *httptest.Server
	tstAuthRequest *entity.AuthRequest = nil
	tstAuthorizationCode = "abcdefghij9876543210"
)

const tstDefaultConfigFile = "../../test/resources/config-acceptancetests.yaml"

func tstSetup(configFilePath string) {
	tstSetupConfig(configFilePath)
	tstSetupHttpTestServer()
	database.Open()

	tstAuthRequest = &entity.AuthRequest{
		Application:      "example-service",
		State:            "Km9NNMK2mx903nlcfkjHd39cdh",
		PkceCodeVerifier: "Nbk2bKbd3klbkkiNKG2cv093hklHKMIHOLKHJacfwklm30m9ym23oHHGGFDSHu9",
		DropOffUrl:       "https://example.com/drop_off_url?dingbaz=5",
		ExpiresAt:        time.Now().Add(config.AuthRequestTimeout()),
	}

	database.GetRepository().AddAuthRequest(context.TODO(), tstAuthRequest)

	dropoffctl.IDPClient = &mockIDPClient{}
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
