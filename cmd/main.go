package main

import (
	"github.com/eurofurence/reg-auth-service/internal/web/app"
	"os"
)

func main() {
	os.Exit(app.New().Run())
}
