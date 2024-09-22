package main

import (
	"fmt"
	"log"

	controller "github.com/Orken1119/Ozinshe/internal/controllers"
	pkg "github.com/Orken1119/Ozinshe/pkg"
	"github.com/gin-gonic/gin"
)

// @title Ozinshe API
// @version 1.0
// @description This is the API server for Ozinshe application.
// @termsOfService http://your-terms-of-service-url.com

// @contact.name API Support
// @contact.url http://your-support-url.com
// @contact.email support@ozinshe.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:1136
// @BasePath /

// @securityDefinitions.basic BasicAuth
func main() {
	app, err := pkg.App()

	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDBConnection()

	ginRouter := gin.Default()

	controller.Setup(app, ginRouter)

	ginRouter.Run(fmt.Sprintf(":%d", 1136))
}
